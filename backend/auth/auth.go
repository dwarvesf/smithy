package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth"
)

const (
	// admin will be able to do everything (get, post, put, delete, ..)
	Admin string = "admin"
	// user can only get data
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
	})

	if err != nil {
		return ""
	}

	return tokenString
}

// Verifier use for verify jwt
func (jwt *JWT) VerifierHandler() func(http.Handler) http.Handler {
	return jwtauth.Verifier(jwt.TokenAuth)
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, claims, _ := jwtauth.FromContext(r.Context())

		uri := r.RequestURI[:6]

		if uri == "/query" {
			if claims["role"] != Admin && claims["role"] != User {
				http.Error(w, http.StatusText(401), 401)
				return
			}
		} else {
			if claims["role"] != Admin {
				http.Error(w, http.StatusText(401), 401)
				return
			}
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
