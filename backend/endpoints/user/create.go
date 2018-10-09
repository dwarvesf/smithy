package user

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// CreateRequest request for add view
type CreateRequest struct {
	User domain.User `json:"user"`
}

// CreateResponse response for add view
type CreateResponse struct {
	Status string      `json:"status"`
	User   domain.User `json:"user"`
}

// MakeCreateUserEndpoint .
func MakeCreateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		err := s.UserService.Create(&req.User)
		if err != nil {
			return nil, err
		}

		return CreateResponse{Status: "success", User: req.User}, nil
	}
}
