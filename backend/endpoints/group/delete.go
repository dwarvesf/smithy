package group

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// DeleteRequest request for add view
type DeleteRequest struct {
	GroupID domain.UUID `json:"-"`
}

// DeleteResponse response for add view
type DeleteResponse struct {
	Status string `json:"status"`
}

// MakeDeleteGroupEndpoint .
func MakeDeleteGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DeleteRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		group := &domain.Group{Model: domain.Model{ID: req.GroupID}}
		err := s.GroupService.Delete(group)
		if err != nil {
			return nil, err
		}

		return DeleteResponse{Status: "success"}, nil
	}
}
