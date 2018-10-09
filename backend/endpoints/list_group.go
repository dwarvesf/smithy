package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// ListGroupResponse response for list view
type ListGroupResponse struct {
	Status string                `json:"status"`
	Groups []backendConfig.Group `json:"groups"`
}

func makeListGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return ListGroupResponse{
			Status: "success",
			Groups: s.Wrapper.SyncConfig().Authentication.Groups,
		}, nil
	}
}
