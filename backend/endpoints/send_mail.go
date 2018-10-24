package endpoints

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"strconv"

	"github.com/dwarvesf/smithy/backend/domain"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/email"
	"github.com/dwarvesf/smithy/backend/service"
)

// SendMailRequest store SendMail structer
type SendEmailRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// SendMailResponse store SendMail respone
type SendEmailResponse struct {
	Status string `json:"status"`
}

func makeSendEmailEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(SendEmailRequest)

		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		var (
			userName  = req.Username
			userEmail = req.Email
		)

		//generate confirm code
		code := strconv.Itoa(rand.Intn(899999) + 100000)

		sendgrid := email.New(os.Getenv("SENDGRID_API_KEY"))
		if err := sendgrid.Send(userName, userEmail, code); err != nil {
			return nil, err
		}

		_, err := s.UserService.Update(&domain.User{Username: userName, ConfirmCode: code})
		if err != nil {
			return nil, err
		}

		return SendEmailResponse{"success"}, nil
	}
}
