package endpoints

import (
	"context"
	"errors"
	"os"

	"github.com/go-kit/kit/endpoint"

	auth "github.com/dwarvesf/smithy/backend/auth"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// ConnectGoogleRequest
type ConnectGoogleRequest struct {
	Code string `json:"code"`
}

// ConnectGoogleResponse
type ConnectGoogleResponse struct {
	AccessToken string `json:"access_token"`
}

func makeConnectGoogleEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ConnectGoogleRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		email, err := s.Provider.CompleteUserAuth(req.Code)
		if err != nil {
			return nil, err
		}

		user := &backendConfig.User{Email: email, IsEmailAccount: true, Role: "user"}

		cfg := s.Wrapper.SyncConfig()
		if cfg.Authentication.IsEmailAccountExist(email.ID) {
			cfg.Authentication.UpdateUser(user)
		} else {
			cfg.Authentication.AddUser(user)
		}

		wr := backendConfig.WriteYAML(os.Getenv("CONFIG_FILE_PATH"))
		if err := wr.Write(cfg); err != nil {
			return nil, err
		}

		// create token jwt
		authToken := auth.NewAuthenticate(cfg, auth.SetUserID(email.ID), auth.SetRole("user"), auth.SetIsEmailAccount(true))

		return ConnectGoogleResponse{authToken.Encode()}, nil
	}
}
