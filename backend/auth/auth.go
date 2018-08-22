package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/k0kubun/pp"
)

//JWT for user authenticaion
type JWT struct {
	Username  string
	Rule      string
	TokenAuth *jwtauth.JWTAuth
}

func New(secretKey string, username string, rule string) *JWT {
	jwt := &JWT{
		Username:  username,
		Rule:      rule,
		TokenAuth: jwtauth.New("HS256", []byte(secretKey), nil),
	}

	return jwt
}

func (jwt *JWT) Encode(username string, rule string) string {
	claims := jwtauth.Claims{}.
		Set(jwt.Username, username).
		Set(jwt.Rule, rule).
		SetIssuedNow()
	_, tokenString, err := jwt.TokenAuth.Encode(claims)

	if err != nil {
		pp.Println(err)
		return ""
	}

	return tokenString
}

func (jwt *JWT) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(jwt.TokenAuth)
}
