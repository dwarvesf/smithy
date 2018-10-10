package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/common/database"
)

const (
	//Admin will be able to do everything (get, post, put, delete, ..)
	Admin string = "admin"
	//User can only get data
	User string = "user"
)

//JWT for user authenticaion
type JWT struct {
	Username  string
	Role      string
	TokenAuth *jwtauth.JWTAuth
}

// New use in backend.go, use for create jwt object
func New(secretKey string, username string, role string) *JWT {
	{
	}
	jwt := &JWT{
		Username:  username,
		Role:      role,
		TokenAuth: jwtauth.New("HS256", []byte(secretKey), nil),
	}

	return jwt
}

// Encode use for encode jwt
func (jwt *JWT) Encode() string {
	_, tokenString, err := jwt.TokenAuth.Encode(jwtauth.Claims{
		"username": jwt.Username,
		"role":     jwt.Role,
	}.SetExpiryIn(time.Second * 3600 * 100))

	if err != nil {
		return ""
	}

	return tokenString
}

//VerifierHandler use for verify jwt
func (jwt *JWT) VerifierHandler() func(http.Handler) http.Handler {
	return jwtauth.Verifier(jwt.TokenAuth)
}

// URIType .
type URIType int

const (
	URITypeAgentSync URIType = 1
	URITypeCRUD      URIType = 2
	URITypeGroup     URIType = 3
)

// parse uri => type, dbName, tableName, method, ok
func parseURI(uri string) (URIType, string, string, string, bool) {
	uriParts := strings.Split(uri, "/")
	if len(uriParts) <= 0 {
		return 0, "", "", "", false
	}

	if uriParts[1] == "groups" {
		return URITypeGroup, "", "", "", true
	}

	if len(uriParts) <= 5 {
		return URITypeAgentSync, "", "", "", true
	}

	// dbName, tableName, method
	return URITypeCRUD, uriParts[2], uriParts[4], uriParts[5], true
}

//Authorization return json in middleware authorization
func Authorization(cfg *backendConfig.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			userName := claims["username"].(string)

			// sample uri: /databases/fortress/table/users/create
			uriType, dbName, tableName, method, ok := parseURI(r.RequestURI)

			if !ok {
				encodeJSONError(ErrInvalidURL, w)
				return
			} else if uriType == URITypeAgentSync || uriType == URITypeGroup {
				// case /agent-sync
				if claims["role"] != Admin {
					encodeJSONError(ErrUnauthorized, w)
					return
				}
			} else if uriType == URITypeCRUD {
				if claims["role"] != Admin && claims["role"] != User {
					encodeJSONError(ErrUnauthorized, w)
					return
				}

				// check dbName is invalid in agent config
				model, ok := cfg.ModelMap[dbName]
				if !ok {
					encodeJSONError(ErrInvalidDatabaseName, w)
					return
				}

				// check table name is invalid in agent config
				tableInfo, ok := model[tableName]
				if !ok {
					encodeJSONError(ErrInvalidTableName, w)
					return
				}

				// get permission (user && group)
				permissions, err := cfg.Authentication.GetFinalPermission(userName)
				if err != nil {
					encodeJSONError(err, w)
					return
				}

				// default permission
				finalPermission := &database.ACLDetail{
					Insert: false,
					Delete: false,
					Select: false,
					Update: false,
				}
				_, ok = permissions[dbName]
				if ok {
					p, ok := permissions[dbName][tableName]
					if ok {
						finalPermission = p
					}
				}

				// user just can access the url when user has user permisstion or table permisstion
				if err := authorizeCRUD(method, finalPermission, tableInfo.ACL); err != nil {
					encodeJSONError(err, w)
					return
				}
			}
			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}

//Authenticator use for authentication user
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			encodeJSONError(err, w)
			return
		}

		if token == nil || !token.Valid {
			encodeJSONError(err, w)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func encodeJSONError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	// enforce json response
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func authorizeCRUD(method string, acl *database.ACLDetail, ACLTable string) error {
	// if user hadn't user permisstion or table permisstion. They would be rejected
	switch method {
	case "query":
		if !acl.Select || !strings.ContainsAny(ACLTable, "r") {
			return ErrUnauthorized
		}
	case "create":
		if !acl.Insert || !strings.ContainsAny(ACLTable, "c") {
			return ErrUnauthorized
		}
	case "update":
		if !acl.Update || !strings.ContainsAny(ACLTable, "u") {
			return ErrUnauthorized
		}
	case "delete":
		if !acl.Delete || !strings.ContainsAny(ACLTable, "d") {
			return ErrUnauthorized
		}
	default:
		return ErrInvalidHTTPMethod
	}
	return nil
}
