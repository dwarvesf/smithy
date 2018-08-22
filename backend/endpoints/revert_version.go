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
	VersionNumber int64 `json:"version"`
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

		cfg := s.Config.Config()
		reader := config.NewBoltPersistent(cfg.PersistenceDB, req.VersionNumber)
		vCfg, err := reader.Read()

		if err != nil {
			return nil, err
		}

		if err := s.Config.Config().UpdateConfig(vCfg); err != nil {
			return nil, err
		}

		return RevertVersionResponse{s.Config.Config().Version}, nil
	}
}
