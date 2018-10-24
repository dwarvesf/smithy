package endpoints

import (
	"context"
	"errors"

	"github.com/dwarvesf/smithy/backend/domain"

	"github.com/go-kit/kit/endpoint"

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

		user, err := s.UserService.Find(&domain.User{Username: req.Username})
		if err != nil {
			return nil, err
		}

		return FindAccountResponse{user.Email}, nil
	}
}
