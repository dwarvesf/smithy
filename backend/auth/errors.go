package auth

import "net/http"

var (
	ErrLogin = errLogin{}
)

type errLogin struct{}

func (errLogin) Error() string {
	return "User name and password is invalid"
}

func (errLogin) StatusCode() int {
	return http.StatusUnauthorized
}
