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
)

// AvailableModelsResponse response for available model endpoints
type AvailableModelsResponse struct {
	Status           string           `json:"status"`
	AvailableMethods []string         `json:"available_methods"`
	Models           []database.Model `json:"models"`
}

func makeAvailableModelsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data := s.Config.Config().ModelList

		return AvailableModelsResponse{"success", availableMethods, data}, nil
	}
}
