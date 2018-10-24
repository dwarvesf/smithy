package group

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// FindRequest request for add view
type FindRequest struct {
	GroupID domain.UUID `json:"-"`
}

// FindResponse response for list view
type FindResponse struct {
	Status string        `json:"status"`
	Group  *domain.Group `json:"group"`
}

// MakeGroupFindEndpoint .
func MakeGroupFindEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FindRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		group := &domain.Group{Model: domain.Model{ID: req.GroupID}}
		group, err := s.GroupService.Find(group)

		if err != nil {
			return nil, err
		}

		return FindResponse{
			Status: "success",
			Group:  group,
		}, nil
	}
}
