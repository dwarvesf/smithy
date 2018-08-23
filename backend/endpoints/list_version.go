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

		querier := config.NewBoltIO(cfg.PersistenceDB, 0)
		versions := querier.ListVersion()

		return listVersionResponse{versions}, nil
	}
}
