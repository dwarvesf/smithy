package auth

import "net/http"

var (
	//ErrLogin return status 401 in login error
	ErrLogin = errLogin{}
)

type errLogin struct{}

func (errLogin) Error() string {
	return "User name and password is invalid"
}

func (errLogin) StatusCode() int {
	return http.StatusUnauthorized
}

//ErrAuthentication use for middleware in authentication
type ErrAuthentication struct {
	Error string `json:"error"`
}
