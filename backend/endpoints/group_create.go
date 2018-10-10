package endpoints

import (
	"context"
	"errors"
	"os"

	"github.com/go-kit/kit/endpoint"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// CreateGroupRequest request for add view
type CreateGroupRequest struct {
	Group backendConfig.Group `json:"group"`
}

// CreateGroupResponse response for add view
type CreateGroupResponse struct {
	Status string              `json:"status"`
	Group  backendConfig.Group `json:"group"`
}

func makeCreateGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateGroupRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		cfg := s.Wrapper.SyncConfig()
		cfg.Authentication.AddGroup(&req.Group)

		wr := backendConfig.WriteYAML(os.Getenv("CONFIG_FILE_PATH"))
		if err := wr.Write(cfg); err != nil {
			return nil, err
		}

		return CreateGroupResponse{Status: "success", Group: req.Group}, nil
	}
}
