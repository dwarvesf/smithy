package group

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// CreateRequest request for add view
type CreateRequest struct {
	Group domain.Group `json:"group"`
}

// CreateResponse response for add view
type CreateResponse struct {
	Status string       `json:"status"`
	Group  domain.Group `json:"group"`
}

// MakeCreateGroupEndpoint .
func MakeCreateGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		err := s.GroupService.Create(&req.Group)
		if err != nil {
			return nil, err
		}

		return CreateResponse{Status: "success", Group: req.Group}, nil
	}
}
