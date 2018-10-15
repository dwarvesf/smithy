package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	jwtAuth "github.com/dwarvesf/smithy/backend/auth"
	BackendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// FindAccountRequest store login structer
type FindAccountRequest struct {
	Username string `json:"username"`
}

// FindAccountResponse store login respone
type FindAccountResponse struct {
	Email string `json:"email"`
}

func makeFindAccountEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FindAccountRequest)

		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		email, err := confirm(req.Username, s.SyncConfig().ConvertUserListToMap())
		if err != nil {
			return nil, err
		}

		return FindAccountResponse{email}, nil
	}
}

func confirm(username string, users map[string]BackendConfig.User) (string, error) {
	userInfo, ok := users[username]
	if !ok {
		return "", jwtAuth.ErrUserNameIsNotExist
	}
	return userInfo.Email, nil
}
