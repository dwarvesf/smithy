package endpoints

import (
	"context"
	"os"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
)

// AuthInformationResponse store auth information
type AuthInformationResponse struct {
	ApplicationName string   `json:"application_name"`
	ClientID        string   `json:"client_id"`
	ClientSecret    string   `json:"client_secret"`
	Scopes          []string `json:"scopes"`
}

func makeAuthInformationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return AuthInformationResponse{
			ApplicationName: "smithy",
			ClientID:        os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret:    os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes: []string{
				"profile",
				"email",
				"https://www.googleapis.com/auth/plus.login"},
		}, nil
	}
}
