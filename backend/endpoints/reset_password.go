package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/auth"
	"github.com/dwarvesf/smithy/backend/domain"
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
			return nil, auth.ErrRePasswordIsNotMatch
		}

		complexity := checkPassword(newPasswordConfirmation)
		if complexity == VeryWeak || complexity == TooShort {
			return nil, auth.ErrPassWordIsVeryWeak
		}

		_, err := s.UserService.Update(&domain.User{Username: userName, Password: newPassword})
		if err != nil {
			return nil, err
		}

		return ResetPasswordResponse{complexity}, nil
	}
}
