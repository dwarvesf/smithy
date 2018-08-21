package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/k0kubun/pp"

	"github.com/dwarvesf/smithy/backend/jwt"
	"github.com/dwarvesf/smithy/backend/service"
)

// LoginRequest store login structer
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse store login respone
type LoginResponse struct {
	Authentication string `json:"authentication"`
}

func makeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		loginAuth := jwt.JWT{}.New(s.SecretKey)

		req, ok := request.(LoginRequest)

		if !ok {
			return nil, errors.New("Login fail")
		}

		// check username and password is existed in persistant,
		// if exist, will return jwt with role and username
		// otherwise return login fail

		// dummy rule
		rule := "admin"

		// two line code below to debug
		pp.Println(req.Username)
		pp.Println(req.Password)

		// if login fail
		if !login(req.Username, req.Password) {
			return nil, errors.New("Login fail")
		}

		// login success
		// return json with jwt is attached

		return LoginResponse{loginAuth.Encode(req.Username, rule)}, nil
	}
}

func login(username, password string) bool {
	return true
}
