package endpoints

import (
	"context"
	"errors"
	"log"
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

		cfg := s.SyncConfig()
		userMap := cfg.ConvertUserListToMap()
		userInfo, ok := userMap[userName]
		if !ok {
			return nil, errors.New("username is invalid")
		}

		if newPassword != newPasswordConfirmation {
			return nil, jwtAuth.ErrRePasswordIsNotMatch
		}

		complexity := checkPassword(newPasswordConfirmation)
		if complexity == VeryWeak || complexity == TooShort {
			return nil, jwtAuth.ErrPassWordIsVeryWeak
		}

		tmpCfg := cloneConfig(cfg)
		updatePassword(tmpCfg, cfg, userInfo, newPassword)

		// config name
		wr := backendConfig.WriteYAML(os.Getenv("CONFIG_FILE_PATH"))
		if err := wr.Write(tmpCfg); err != nil {
			log.Fatalln(err)
		}

		return ResetPasswordResponse{complexity}, nil
	}
}
