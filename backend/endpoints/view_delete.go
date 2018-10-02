package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
)

// DeleteViewRequest request for add view
type DeleteViewRequest struct {
	DatabaseName string `json:"-"`
	SQLID        int    `json:"-"`
}

// DeleteViewResponse response for add view
type DeleteViewResponse struct {
	Status string `json:"status"`
}

func makeDeleteViewEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DeleteViewRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		err := s.WriteReadDeleter.Delete(req.SQLID)
		if err != nil {
			return nil, err
		}

		return DeleteViewResponse{Status: "success"}, nil
	}
}
