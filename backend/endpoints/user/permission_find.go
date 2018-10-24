package user

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// PermissionFindRequest request for add view
type PermissionFindRequest struct {
	UserID       domain.UUID `json:"-"`
	TableName    string      `json:"table_name"`
	DatabaseName string      `json:"database_name"`
}

// PermissionFindResponse response for list view
type PermissionFindResponse struct {
	Status      string              `json:"status"`
	Permissions []domain.Permission `json:"permissions"`
}

// MakeUserPermissionFindEndpoint .
func MakeUserPermissionFindEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(PermissionFindRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		user := &domain.User{Model: domain.Model{ID: req.UserID}}

		var (
			permissions []domain.Permission
			err         error
		)
		if req.TableName != "" && req.DatabaseName != "" {
			permissions, err = s.UserService.GetPermission(user, req.DatabaseName, req.TableName)
		} else {
			permissions, err = s.PermissionService.FindByUser(user)
		}

		if err != nil {
			return nil, err
		}

		return PermissionFindResponse{
			Status:      "success",
			Permissions: permissions,
		}, nil
	}
}
