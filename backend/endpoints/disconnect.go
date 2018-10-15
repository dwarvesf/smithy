package endpoints

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
)

// DisConnectGoogleResponse
type DisConnectGoogleResponse struct {
	Status string `json:"status"`
}

func makeDisConnectGoogleEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, claims, _ := jwtauth.FromContext(ctx)
		var (
			emailID = claims["email_id"].(string)
		)

		user := s.SyncConfig().Authentication.FindUserByEmailID(emailID)
		if user == nil {
			return nil, errors.New("Current user not connected")
		}

		// Execute HTTP GET request to revoke current token
		url := "https://accounts.google.com/o/oauth2/revoke?token=" + user.Email.AccessToken
		resp, err := http.Get(url)
		if err != nil {
			return nil, errors.New("Failed to revoke token for a given user")
		}
		defer resp.Body.Close()

		return DisConnectGoogleResponse{"success"}, nil
	}
}
