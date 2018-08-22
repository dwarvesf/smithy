package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

type revertVersionResquest struct {
	Version string `json:"version"`
}

type revertVersionResponse struct {
	Status string `json:"status"`
}

func makeRevertVersionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(revertVersionResquest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		cfg := s.Config.Config()
		reader := config.NewBoltReader(req.Version, cfg.PersistenceDB)
		vCfg, err := reader.Read()

		if err != nil {
			return nil, err
		}

		if err := s.Config.Config().UpdateConfig(vCfg); err != nil {
			return nil, err
		}

		return revertVersionResponse{"success"}, nil
	}
}
