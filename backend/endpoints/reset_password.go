package endpoints

import (
	"context"
	"errors"
	"os"

	"github.com/go-kit/kit/endpoint"

	jwtAuth "github.com/dwarvesf/smithy/backend/auth"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/service"
)

// ResetPasswordRequest store reset password structer
type ResetPasswordRequest struct {
	Username                string `json:"username"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

// ResetPasswordResponse store reset password respone
type ResetPasswordResponse struct {
	Complexity string `json:"complexity"`
}

func makeResetPasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ResetPasswordRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		var (
			userName                = req.Username
			newPassword             = req.NewPassword
			newPasswordConfirmation = req.NewPasswordConfirmation
		)

		if newPassword != newPasswordConfirmation {
			return nil, jwtAuth.ErrRePasswordIsNotMatch
		}

		complexity := checkPassword(newPasswordConfirmation)
		if complexity == VeryWeak || complexity == TooShort {
			return nil, jwtAuth.ErrPassWordIsVeryWeak
		}

		cfg := s.SyncConfig()
		cfg.Authentication.UpdateConfirmCode(userName, newPassword)

		// config name
		wr := backendConfig.WriteYAML(os.Getenv("CONFIG_FILE_PATH"))
		if err := wr.Write(cfg); err != nil {
			return nil, err
		}

		return ResetPasswordResponse{complexity}, nil
	}
}
