package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

type listVersionResponse struct {
	Versions []config.Version `json:"versions"`
}

func makeListVersionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cfg := s.Config.Config()

		querier := config.NewBoltPersistent(cfg.PersistenceFileName, 0)
		versions, err := querier.ListVersion()

		if err != nil {
			return nil, err
		}

		return listVersionResponse{versions}, nil
	}
}
