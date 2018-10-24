package group

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// FindAllResponse response for list view
type FindAllResponse struct {
	Status string         `json:"status"`
	Groups []domain.Group `json:"groups"`
}

// MakeGroupFindAllEndpoint .
func MakeGroupFindAllEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		groups, err := s.GroupService.FindAll()
		if err != nil {
			return nil, err
		}

		return FindAllResponse{
			Status: "success",
			Groups: groups,
		}, nil
	}
}
