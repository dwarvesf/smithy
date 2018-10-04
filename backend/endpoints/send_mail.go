package endpoints

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/go-kit/kit/endpoint"

	jwtAuth "github.com/dwarvesf/smithy/backend/auth"
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

const (
	YOUR_SENDGRID_API_KEY = "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"
)

func makeSendEmailEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(SendEmailRequest)

		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		cfg := s.SyncConfig()
		userMap := cfg.ConvertUserListToMap()
		userInfo, ok := userMap[req.Username]
		if !ok {
			return nil, jwtAuth.ErrUserNameIsNotExist
		}
		//generate confirm code
		code := strconv.Itoa(rand.Intn(899999) + 100000)

		sendgrid := email.New(YOUR_SENDGRID_API_KEY)
		if err := sendgrid.Send(req.Username, req.Email, code); err != nil {
			return nil, err
		}

		tmpCfg := cloneConfig(cfg)
		updateConfirmCode(tmpCfg, cfg, userInfo, code)

		// config name
		wr := backendConfig.WriteYAML(os.Getenv("CONFIG_FILE_PATH"))
		if err := wr.Write(tmpCfg); err != nil {
			log.Fatalln(err)
		}

		return SendEmailResponse{"success"}, nil
	}
}

func updateConfirmCode(tmpCfg, cfg *backendConfig.Config, userInfo backendConfig.User, newCode string) {
	for i, user := range cfg.Authentication.UserList {
		if user.Username == userInfo.Username {
			cfg.Authentication.UserList[i].ConfirmCode = newCode
		}
	}
	tmpCfg.Authentication = cfg.Authentication
}
