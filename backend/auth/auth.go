package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
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

//Authorization return json in middleware authorization
func Authorization(cfg *backendConfig.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			modelMap, userMap := cfg.ModelMap, cfg.ConvertUserListToMap()
			userInfo, ok := userMap[claims["username"].(string)]
			if !ok {
				encodeJSONError(ErrInvalidUserName, w)
				return
			}
			// POST: localhost:2999/databases/fortress/users/create
			uriParts := strings.Split(r.RequestURI, "/")
			if len(uriParts) <= 0 {
				encodeJSONError(ErrInvalidURL, w)
				return
			} else if len(uriParts) <= 3 {
				// case /agent-sync
				if claims["role"] != Admin {
					encodeJSONError(ErrUnauthorized, w)
					return
				}
			} else if len(uriParts) <= 5 {
				// case /agent-sync
				if claims["role"] != Admin {
					encodeJSONError(ErrUnauthorized, w)
					return
				}
			} else {
				if claims["role"] != Admin && claims["role"] != User {
					encodeJSONError(ErrUnauthorized, w)
					return
				}
				var (
					//fortress
					dbName = uriParts[2]
					//users
					tableName = uriParts[4]
				)
				existDB := false
				for _, db := range userInfo.DatabaseList {
					// check dbName in URL with dbName in dashboard config
					if db.DBName == dbName {
						existDB = true
						// check dbName is invalid in agent config
						model, ok := modelMap[dbName]
						if !ok {
							encodeJSONError(ErrInvalidDatabaseName, w)
							return
						}
						existTable := true
						for _, table := range db.Tables {
							if table.TableName == tableName {
								existTable = true
								// check table name is invalid in agent config
								tableInfo, ok := model[tableName]
								if !ok {
									encodeJSONError(ErrInvalidTableName, w)
									return
								}
								// get table ACL in agent config
								ACLTable := tableInfo.ACL
								// set ACLDeltail
								table.MakeACLDetailFromACL()
								// user just can access the url when user has user permisstion or table permisstion
								if err := authorizeCRUD(r.Method, table, ACLTable); err != nil {
									encodeJSONError(err, w)
									return
								}
							}
						}
						if !existTable {
							encodeJSONError(ErrInvalidTableName, w)
							return
						}
					}
				}
				if !existDB {
					encodeJSONError(ErrInvalidDatabaseName, w)
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

func authorizeCRUD(method string, table backendConfig.Table, ACLTable string) error {
	// if user hadn't user permisstion or table permisstion. They would be rejected
	switch method {
	case "GET":
		if !table.ACLDetail.Select || !strings.ContainsAny(ACLTable, "r") {
			return ErrUnauthorized
		}
	case "POST":
		if !table.ACLDetail.Insert || !strings.ContainsAny(ACLTable, "c") {
			return ErrUnauthorized
		}
	case "PUT":
		if !table.ACLDetail.Update || !strings.ContainsAny(ACLTable, "u") {
			return ErrUnauthorized
		}
	case "DELETE":
		if !table.ACLDetail.Delete || !strings.ContainsAny(ACLTable, "d") {
			return ErrUnauthorized
		}
	default:
		return ErrInvalidHTTPMethod
	}
	return nil
}
