package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// FindAllResponse response for list view
type FindAllResponse struct {
	Status string        `json:"status"`
	Users  []domain.User `json:"users"`
}

// MakeUserFindAllEndpoint .
func MakeUserFindAllEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.UserService.FindAll()
		if err != nil {
			return nil, err
		}

		return FindAllResponse{
			Status: "success",
			Users:  users,
		}, nil
	}
}
