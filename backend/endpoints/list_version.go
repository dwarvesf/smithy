package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/config/persistence"
	"github.com/dwarvesf/smithy/backend/service"
)

type listVersionResponse struct {
	Versions []string `json:"versions"`
}

func makeListVersionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cfg := s.Config.Config()

		pio := persistence.NewBoltPersistence(cfg.PersistenceFileName, cfg.PersistenceDB)
		versions := pio.ListVersion()

		return listVersionResponse{versions}, nil
	}
}
