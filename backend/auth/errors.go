package auth

import "net/http"

var (
	//ErrLogin return status 401 in login error
	ErrLogin                = errLogin{}
	ErrUnauthorized         = errUnauthorized{}
	ErrInvalidUserName      = errInvalidUserName{}
	ErrInvalidHTTPMethod    = errInvalidHTTPMethod{}
	ErrInvalidURL           = errInvalidURL{}
	ErrInvalidDatabaseName  = errInvalidDatabaseName{}
	ErrInvalidTableName     = errInvalidTableName{}
	ErrOldPasswordInvalid   = errOldPasswordInvalid{}
	ErrRePasswordIsNotMatch = errRePasswordIsNotMatch{}
	ErrPassWordIsVeryWeak   = errPassWordIsVeryWeak{}
)

type errLogin struct{}

func (errLogin) Error() string {
	return "user name or password is invalid"
}

func (errLogin) StatusCode() int {
	return http.StatusUnauthorized
}

type errUnauthorized struct{}

func (errUnauthorized) Error() string {
	return "Unauthorized"
}

func (errUnauthorized) StatusCode() int {
	return http.StatusUnauthorized
}

type errInvalidUserName struct{}

func (errInvalidUserName) Error() string {
	return "Username is invalid"
}

func (errInvalidUserName) StatusCode() int {
	return http.StatusUnauthorized
}

type errInvalidHTTPMethod struct{}

func (errInvalidHTTPMethod) Error() string {
	return "Invalid HTTP method"
}

func (errInvalidHTTPMethod) StatusCode() int {
	return http.StatusUnauthorized
}

type errInvalidURL struct{}

func (errInvalidURL) Error() string {
	return "Invalid URL"
}

func (errInvalidURL) StatusCode() int {
	return http.StatusUnauthorized
}

type errInvalidDatabaseName struct{}

func (errInvalidDatabaseName) Error() string {
	return "Unknown database name"
}

func (errInvalidDatabaseName) StatusCode() int {
	return http.StatusUnauthorized
}

type errInvalidTableName struct{}

func (errInvalidTableName) Error() string {
	return "Unknown table name"
}

func (errInvalidTableName) StatusCode() int {
	return http.StatusUnauthorized
}

type errOldPasswordInvalid struct{}

func (errOldPasswordInvalid) Error() string {
	return "old password is not correct"
}

func (errOldPasswordInvalid) StatusCode() int {
	return http.StatusUnauthorized
}

type errRePasswordIsNotMatch struct{}

func (errRePasswordIsNotMatch) Error() string {
	return "new password confirmation is not match"
}

func (errRePasswordIsNotMatch) StatusCode() int {
	return http.StatusUnauthorized
}

type errPassWordIsVeryWeak struct{}

func (errPassWordIsVeryWeak) Error() string {
	return "new password is very weak! Please try again"
}

func (errPassWordIsVeryWeak) StatusCode() int {
	return http.StatusUnauthorized
}

//ErrAuthentication use for middleware in authentication
type ErrAuthentication struct {
	Error string `json:"error"`
}
