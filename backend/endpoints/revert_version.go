package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// RevertVersionResquest request for revert version
type RevertVersionResquest struct {
	VersionID int `json:"version_id"`
}

// RevertVersionResponse response for revert version
type RevertVersionResponse struct {
	Version config.Version `json:"version"`
}

func makeRevertVersionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(RevertVersionResquest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		if err := s.Config.Config().ChangeVersion(req.VersionID); err != nil {
			return nil, err
		}

		return RevertVersionResponse{s.Config.Config().Version}, nil
	}
}
