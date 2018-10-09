package endpoints

import (
	"context"
	"errors"
	"os"

	"github.com/go-kit/kit/endpoint"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// DeleteGroupRequest request for add view
type DeleteGroupRequest struct {
	GroupID string `json:"-"`
}

// DeleteGroupResponse response for add view
type DeleteGroupResponse struct {
	Status string `json:"status"`
}

func makeDeleteGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DeleteGroupRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		cfg := s.Wrapper.SyncConfig()
		cfg.Authentication.DeleteGroup(req.GroupID)

		wr := backendConfig.WriteYAML(os.Getenv("CONFIG_FILE_PATH"))
		if err := wr.Write(cfg); err != nil {
			return nil, err
		}

		return DeleteGroupResponse{Status: "success"}, nil
	}
}
