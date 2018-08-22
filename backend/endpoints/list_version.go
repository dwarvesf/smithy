package endpoints

import (
	"context"

	"github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
	"github.com/go-kit/kit/endpoint"
)

type listVersionResponse struct {
	Versions []string `json:"versions"`
}

func makeListVersionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cfg := s.Config.Config()

		querier := config.NewBoltQuerier(cfg.PersistenceDB)
		versions := querier.ListVersion()

		return listVersionResponse{versions}, nil
	}
}
