package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/backend/view"
)

// ListViewRequest request for list view
type ListViewRequest struct {
	DatabaseName string `json:"-"`
}

// ListViewResponse response for list view
type ListViewResponse struct {
	Status string       `json:"status"`
	Views  []*view.View `json:"views"`
}

func makeListViewEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ListViewRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		sqlcmds, err := s.WriterReaderDeleter.ListCommandsByDBName(req.DatabaseName)
		if err != nil {
			return nil, err
		}

		return ListViewResponse{Status: "success", Views: sqlcmds}, nil
	}
}
