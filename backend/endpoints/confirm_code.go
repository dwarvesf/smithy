package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	jwtAuth "github.com/dwarvesf/smithy/backend/auth"
	BackendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// ConfirmCodeRequest store Confirm code structer
type ConfirmCodeRequest struct {
	Username    string `json:"username"`
	ConfirmCode string `json:"confirm_code"`
}

// ConfirmCodeResponse store Confirm code respone
type ConfirmCodeResponse struct {
	Status string `json:"status"`
}

func makeConfirmCodeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ConfirmCodeRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		err := confirmCode(req.Username, req.ConfirmCode, s.SyncConfig().ConvertUserListToMap())
		if err != nil {
			return nil, err
		}

		return ConfirmCodeResponse{"success"}, nil
	}
}

func confirmCode(username, confirmCode string, users map[string]BackendConfig.User) error {
	userInfo, ok := users[username]
	if !ok {
		return jwtAuth.ErrUserNameIsNotExist
	}

	if userInfo.ConfirmCode != confirmCode {
		return jwtAuth.ErrConfirmCodeIsNotMatch
	}

	return nil
}
