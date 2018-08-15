package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
)

// Endpoints .
type Endpoints struct {
	AgentSync endpoint.Endpoint

	DBQuery  endpoint.Endpoint
	DBCreate endpoint.Endpoint

	AvailableModels endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		AgentSync: makeAgentSyncEndpoint(s),

		DBQuery:  makeDBQueryEndpoint(s),
		DBCreate: makeDBCreateEndpoint(s),

		AvailableModels: makeAvailableModelsEndpoint(s),
	}
}
