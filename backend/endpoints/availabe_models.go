package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/common/database"
)

var (
	availableMethods = []string{
		"FindByID",
		"FindByColumnName",
		"FindAll",
		"Create",
		"Delete",
		"Update",
	}

	availableHookTypes = database.HookTypes
)

// AvailableModelsResponse response for available model endpoints
type AvailableModelsResponse struct {
	Status             string           `json:"status"`
	AvailableMethods   []string         `json:"available_methods"`
	AvailableHookTypes []string         `json:"available_hook_types"`
	Models             []database.Model `json:"models"`
}

func makeAvailableModelsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data := s.SyncConfig().ModelList

		return AvailableModelsResponse{
			Status:             "success",
			AvailableMethods:   availableMethods,
			Models:             data,
			AvailableHookTypes: availableHookTypes}, nil
	}
}
