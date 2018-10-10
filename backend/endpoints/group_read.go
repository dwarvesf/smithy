package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// ReadGroupRequest request for add view
type ReadGroupRequest struct {
	GroupID string `json:"-"`
}

// ReadGroupResponse response for list view
type ReadGroupResponse struct {
	Status string               `json:"status"`
	Group  *backendConfig.Group `json:"group"`
}

func makeReadGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ReadGroupRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		var group *backendConfig.Group
		cfg := s.Wrapper.SyncConfig()
		groups := cfg.Authentication.Groups
		for i := 0; i < len(groups); i++ {
			if groups[i].ID == req.GroupID {
				group = &groups[i]
				break
			}
		}

		if group == nil {
			return nil, errors.New("group not found")
		}

		return ReadGroupResponse{
			Status: "success",
			Group:  group,
		}, nil
	}
}
