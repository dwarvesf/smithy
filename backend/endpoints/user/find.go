package user

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// FindRequest request for add view
type FindRequest struct {
	UserID domain.UUID `json:"-"`
}

// FindResponse response for list view
type FindResponse struct {
	Status string       `json:"status"`
	User   *domain.User `json:"user"`
}

// MakeUserFindEndpoint .
func MakeUserFindEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FindRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		user := &domain.User{Model: domain.Model{ID: req.UserID}}
		user, err := s.UserService.Find(user)
		if err != nil {
			return nil, err
		}

		return FindResponse{
			Status: "success",
			User:   user,
		}, nil
	}
}
