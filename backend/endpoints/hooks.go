package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// AddHookRequest request for db create data
type AddHookRequest struct {
	TableName   string `json:"table_name"`
	HookContent string `json:"hook_content"`
	HookType    string `json:"hook_type"`
}

// AddHookResponse response for db create data
type AddHookResponse struct {
	Status string `json:"status"`
}

func makeAddHookEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(AddHookRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		cfg := s.SyncConfig()

		err := cfg.AddHook(req.TableName, req.HookType, req.HookContent)
		if err != nil {
			return nil, err
		}

		// Update Config persistent
		writer := config.NewBoltPersistent(cfg.PersistenceFileName, 0)
		err = writer.Write(cfg)
		if err != nil {
			return nil, err
		}

		return AddHookResponse{Status: "success"}, nil
	}
}
