package user

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// UpdateRequest request for add view
type UpdateRequest struct {
	User *domain.User `json:"user"`
}

// UpdateResponse response for list view
type UpdateResponse struct {
	Status string       `json:"status"`
	User   *domain.User `json:"user"`
}

// MakeUpdateUserEndpoint .
func MakeUpdateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdateRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		user, err := s.UserService.Update(req.User)

		if err != nil {
			return nil, err
		}

		return UpdateResponse{
			Status: "success",
			User:   user,
		}, nil
	}
}
