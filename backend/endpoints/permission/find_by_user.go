package permission

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// FindByUserRequest request for add view
type FindByUserRequest struct {
	UserID       domain.UUID `json:"-"`
	TableName    string      `json:"table_name"`
	DatabaseName string      `json:"database_name"`
}

// FindByUserResponse response for list view
type FindByUserResponse struct {
	Status      string              `json:"status"`
	Permissions []domain.Permission `json:"permissions"`
}

// MakeUserPermissionFindEndpoint .
func MakePermissionFindByUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FindByUserRequest)
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

		return FindByUserResponse{
			Status:      "success",
			Permissions: permissions,
		}, nil
	}
}
