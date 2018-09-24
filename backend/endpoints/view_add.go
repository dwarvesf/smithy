package endpoints

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/backend/view"
)

// AddViewRequest request for add view
type AddViewRequest struct {
	view.View
}

// AddViewResponse response for add view
type AddViewResponse struct {
	Status    string      `json:"status"`
	QueryPlan interface{} `json:"query_plan"`
}

func makeAddViewEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(AddViewRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		q, err := req.Validate(s.Mapper)
		if err != nil {
			return nil, err
		}

		req.CreatedAt = time.Now()
		err = s.WriterReaderDeleter.Write(&req.View)
		if err != nil {
			return nil, err
		}

		return AddViewResponse{Status: "success", QueryPlan: q}, nil
	}
}
