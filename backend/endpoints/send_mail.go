package endpoints

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"strconv"

	"github.com/go-kit/kit/endpoint"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
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

		cfg := s.SyncConfig()
		cfg.Authentication.UpdateConfirmCode(userName, code)

		// config name
		wr := backendConfig.WriteYAML(os.Getenv("CONFIG_FILE_PATH"))
		if err := wr.Write(cfg); err != nil {
			return nil, err
		}

		return SendEmailResponse{"success"}, nil
	}
}
