package permission

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// FindByGroupRequest request for add view
type FindByGroupRequest struct {
	GroupID      domain.UUID `json:"-"`
	TableName    string      `json:"table_name"`
	DatabaseName string      `json:"database_name"`
}

// FindByGroupResponse response for list view
type FindByGroupResponse struct {
	Status      string              `json:"status"`
	Permissions []domain.Permission `json:"permissions"`
}

// MakePermissionFindByGroupEndpoint .
func MakePermissionFindByGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FindByGroupRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		group := &domain.Group{Model: domain.Model{ID: req.GroupID}}

		var (
			permissions []domain.Permission
			err         error
		)
		if req.TableName != "" && req.DatabaseName != "" {
			permissions, err = s.GroupService.GetPermission(group, req.DatabaseName, req.TableName)
		} else {
			permissions, err = s.PermissionService.FindByGroup(group)
		}

		if err != nil {
			return nil, err
		}

		return FindByGroupResponse{
			Status:      "success",
			Permissions: permissions,
		}, nil
	}
}
