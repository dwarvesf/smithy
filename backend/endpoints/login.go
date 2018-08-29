package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend"
	jwtAuth "github.com/dwarvesf/smithy/backend/auth"
	BackendConfig "github.com/dwarvesf/smithy/backend/config"
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
		req, ok := request.(LoginRequest)

		if !ok {
			return nil, errors.New("Login fail")
		}

		// check username and password is existed in persistant,
		// if exist, will return jwt with role and username
		// otherwise return login fail

		// if login fail
		ok, rule := login(req.Username, req.Password, s.SyncConfig().ConvertUserListToMap())

		if !ok {
			return nil, jwtAuth.ErrLogin
		}

		// create user authentication
		loginAuth := backend.NewAuthenticate(s.SyncConfig(), req.Username, rule)

		// login success
		// return json with jwt is attached

		return LoginResponse{loginAuth.Encode()}, nil
	}
}

func login(username, password string, users map[string]BackendConfig.User) (bool, string) {
	userInfo, ok := users[username]

	if !ok {
		return false, ""
	}
	if userInfo.Password != password {
		return false, ""
	}

	return true, userInfo.Role
}
