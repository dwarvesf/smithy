package group

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// UpdateRequest request for add view
type UpdateRequest struct {
	Group *domain.Group `json:"group"`
}

// UpdateResponse response for list view
type UpdateResponse struct {
	Status string        `json:"status"`
	Group  *domain.Group `json:"group"`
}

// MakeUpdateGroupEndpoint .
func MakeUpdateGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdateRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		group, err := s.GroupService.Update(req.Group)

		if err != nil {
			return nil, err
		}

		return UpdateResponse{
			Status: "success",
			Group:  group,
		}, nil
	}
}
