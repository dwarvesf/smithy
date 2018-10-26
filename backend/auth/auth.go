package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

const (
	//Admin will be able to do everything (get, post, put, delete, ..)
	Admin string = "admin"
	//User can only get data
	User string = "user"
)

//JWT for user authenticaion
type JWT struct {
	UserName       string
	Email          string
	Role           string
	IsEmailAccount bool
	TokenAuth      *jwtauth.JWTAuth
}

// New use in backend.go, use for create jwt object
func New(secretKey, userName, email, role string, isEmailAccount bool) *JWT {
	jwt := &JWT{
		UserName:       userName,
		Email:          email,
		Role:           role,
		IsEmailAccount: isEmailAccount,
		TokenAuth:      jwtauth.New("HS256", []byte(secretKey), nil),
	}

	return jwt
}

//NewAuthenticate New JWT Authenticain
func NewAuthenticate(c *backendConfig.Config, setters ...Option) *JWT {
	args := &JWT{
		UserName:       "",
		Email:          "",
		Role:           "",
		IsEmailAccount: false,
	}
	for _, setter := range setters {
		setter(args)
	}

	return New(c.Authentication.SerectKey, args.UserName, args.Email, args.Role, args.IsEmailAccount)
}

type Option func(*JWT)

//SetUserID is to set userid
func SetUserName(userName string) Option {
	return func(jwt *JWT) {
		jwt.UserName = userName
	}
}

//SetUserEmail is to set email
func SetEmail(email string) Option {
	return func(jwt *JWT) {
		jwt.Email = email
	}
}

//SetRole is to set role
func SetRole(role string) Option {
	return func(jwt *JWT) {
		jwt.Role = role
	}
}

//SetIsEmailAccount is to set isEmailAccount
func SetIsEmailAccount(isEmailAccount bool) Option {
	return func(jwt *JWT) {
		jwt.IsEmailAccount = isEmailAccount
	}
}

//SetTokenAuth is to set tokenauth
func SetTokenAuth(jwtAuth *jwtauth.JWTAuth) Option {
	return func(jwt *JWT) {
		jwt.TokenAuth = jwtAuth
	}
}

// Encode use for encode jwt
func (jwt *JWT) Encode() string {
	_, tokenString, err := jwt.TokenAuth.Encode(jwtauth.Claims{
		"username":         jwt.UserName,
		"email":            jwt.Email,
		"role":             jwt.Role,
		"is_email_account": jwt.IsEmailAccount,
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
	URITypeAgentSync URIType = iota + 1
	URITypeCRUD
	URITypeGroup
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
func Authorization(cfg *backendConfig.Config, s service.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			user := &domain.User{
				Username:       claims["username"].(string),
				IsEmailAccount: claims["is_email_account"].(bool),
				Email:          claims["email"].(string),
			}

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
				finalPermission, err := s.UserService.GetPermissionUserAndGroup(user, dbName, tableName)
				if err != nil {
					encodeJSONError(err, w)
					return
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

//RequireAdmin return authorization if is admin
func RequireAdmin(s service.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			user := &domain.User{
				Username:       claims["username"].(string),
				IsEmailAccount: claims["is_email_account"].(bool),
				Email:          claims["email"].(string),
				Role:           claims["role"].(string),
			}

			user, err := s.UserService.Find(user)
			if err != nil {
				encodeJSONError(err, w)
				return
			}

			isAdmin := user.Role == Admin
			groups, err := s.GroupService.FindByUser(user)
			if err != nil {
				encodeJSONError(err, w)
				return
			}

			for _, group := range groups {
				isAdmin = isAdmin || (group.Role == Admin)
			}

			if !isAdmin {
				encodeJSONError(errors.New("only admin can access"), w)
				return
			}

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

//Authenticator use for authentication user
func RequireNormalUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		var (
			isEmailAccount = claims["is_email_account"].(bool)
			userName       = claims["username"].(string)
		)

		if userName == "" {
			encodeJSONError(ErrInvalidUserName, w)
			return
		}

		if isEmailAccount {
			encodeJSONError(ErrRequireNormalUser, w)
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

func authorizeCRUD(method string, acl *domain.Permission, ACLTable string) error {
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
