package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
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
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		methods := strings.Split(r.RequestURI, "?")

		var err error
		if len(methods) <= 0 {
			if err = writeToResponse(w, r, "Unauthorized"); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}
		if methods[0] == "/query" {
			if claims["role"] != Admin && claims["role"] != User {
				if err = writeToResponse(w, r, "Unauthorized"); err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
				}
				return
			}
		} else {
			if claims["role"] != Admin {
				if err = writeToResponse(w, r, "Unauthorized"); err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
				}
				return
			}
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

//Authenticator use for authentication user
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			if err = writeToResponse(w, r, err.Error()); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}

		if token == nil || !token.Valid {
			if err = writeToResponse(w, r, err.Error()); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func writeToResponse(w http.ResponseWriter, r *http.Request, errStr string) error {
	w.WriteHeader(http.StatusUnauthorized)
	js, _ := json.Marshal(ErrAuthentication{errStr})
	_, err := w.Write(js)
	return err
}
