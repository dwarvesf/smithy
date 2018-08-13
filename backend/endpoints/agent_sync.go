package endpoints

import (
	"context"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/go-kit/kit/endpoint"
)

type agentSyncResponse struct {
	Status string `json:"status"`
}

func makeAgentSyncEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cfg := s.Config.Config()
		err := cfg.UpdateConfigFromAgent()
		if err != nil {
			return nil, err
		}

		return agentSyncResponse{"success"}, nil
	}
}
