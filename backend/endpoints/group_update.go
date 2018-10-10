package endpoints

import (
	"context"
	"errors"
	"os"

	"github.com/go-kit/kit/endpoint"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// UpdateGroupRequest request for add view
type UpdateGroupRequest struct {
	GroupID string               `json:"-"`
	Group   *backendConfig.Group `json:"group"`
}

// UpdateGroupResponse response for list view
type UpdateGroupResponse struct {
	Status string `json:"status"`
}

func makeUpdateGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdateGroupRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		cfg := s.Wrapper.SyncConfig()
		cfg.Authentication.UpdateGroup(req.GroupID, req.Group)

		wr := backendConfig.WriteYAML(os.Getenv("CONFIG_FILE_PATH"))
		if err := wr.Write(cfg); err != nil {
			return nil, err
		}

		return UpdateGroupResponse{
			Status: "success",
		}, nil
	}
}
