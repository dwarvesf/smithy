package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/auth"
	"github.com/dwarvesf/smithy/backend/service"
)

// ConnectGoogleRequest
type LoginGoogleRequest struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

// ConnectGoogleResponse
type LoginGoogleResponse struct {
	AccessToken string `json:"access_token"`
}

func makeLoginGoogleEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(LoginGoogleRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		user, err := s.Provider.CompleteUserAuth(req.Code, req.RedirectURI)
		if err != nil {
			return nil, err
		}

		if userExist, _ := s.UserService.Find(user); userExist == nil {
			err = s.UserService.Create(user)
		} else {
			user, err = s.UserService.Update(user)
		}

		if err != nil {
			return nil, err
		}

		// create token jwt
		authToken := auth.NewAuthenticate(s.SyncConfig(), auth.SetEmail(user.Email), auth.SetRole(user.Role), auth.SetIsEmailAccount(true))

		return LoginGoogleResponse{authToken.Encode()}, nil
	}
}
