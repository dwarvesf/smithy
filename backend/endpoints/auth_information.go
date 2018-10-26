package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
)

// AuthInformationResponse store auth information
type AuthInformationResponse struct {
	AuthURL string `json:"auth_url"`
}

func makeAuthInformationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return AuthInformationResponse{
			AuthURL: s.Provider.GetAuthURL(),
		}, nil
	}
}
