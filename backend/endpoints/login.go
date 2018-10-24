package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/crypto/bcrypt"

	jwtAuth "github.com/dwarvesf/smithy/backend/auth"
	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// LoginRequest store login structer
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse store login respone
type LoginResponse struct {
	Authentication string `json:"authentication"`
}

func makeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(LoginRequest)

		if !ok {
			return nil, errors.New("Login fail")
		}

		// check username and password is existed in persistant,
		// if exist, will return jwt with role and username
		// otherwise return login fail

		// if login fail
		user := &domain.User{Username: req.Username}
		user, err := s.UserService.Find(user)
		if err != nil {
			return nil, jwtAuth.ErrLogin
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(req.Password))

		if err != nil {
			return nil, jwtAuth.ErrLogin
		}

		// create user authentication
		loginAuth := jwtAuth.NewAuthenticate(s.SyncConfig(), req.Username, user.Role)

		// login success
		// return json with jwt is attached

		return LoginResponse{loginAuth.Encode()}, nil
	}
}
