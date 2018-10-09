package permission

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// UpdateRequest request for add view
type UpdateRequest struct {
	Permission *domain.Permission `json:"permission"`
}

// UpdateResponse response for list view
type UpdateResponse struct {
	Status     string             `json:"status"`
	Permission *domain.Permission `json:"permission"`
}

// MakeUpdatePermissionEndpoint .
func MakeUpdatePermissionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdateRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		permission, err := s.PermissionService.Update(req.Permission)

		if err != nil {
			return nil, err
		}

		return UpdateResponse{
			Status:     "success",
			Permission: permission,
		}, nil
	}
}
