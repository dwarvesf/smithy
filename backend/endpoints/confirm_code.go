package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

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

		cfg := s.SyncConfig()
		userMap := cfg.ConvertUserListToMap()
		userInfo, ok := userMap[req.Username]
		if !ok {
			return nil, errors.New("username is invalid")
		}

		if userInfo.ConfirmCode != req.ConfirmCode {
			return nil, errors.New("confirm code is not match")
		}

		return SendEmailResponse{"success"}, nil
	}
}
