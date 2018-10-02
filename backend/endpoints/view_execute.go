package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/common/database"
)

// ExecuteViewRequest request for add view
type ExecuteViewRequest struct {
	DatabaseName string `json:"-"`
	SQLID        int    `json:"-"`
}

// ExecuteViewResponse response for add view
type ExecuteViewResponse struct {
	Status  string            `json:"status"`
	Columns []string          `json:"columns,omitempty"`
	Rows    []interface{}     `json:"rows,omitempty"`
	Cols    []database.Column `json:"cols,omitempty"`
}

func makeExecuteViewEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ExecuteViewRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		view, err := s.WriteReadDeleter.Read(req.SQLID)
		if err != nil {
			return nil, err
		}

		columns, columnMeta, data, err := s.RawQuery(req.DatabaseName, view.SQL)
		if err != nil {
			return nil, err
		}

		return ExecuteViewResponse{
			Status:  "success",
			Columns: columns,
			Rows:    data,
			Cols:    columnMeta,
		}, nil
	}
}
