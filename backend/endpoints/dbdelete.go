package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
)

// DBDeleteRequest request for db delete data by id
type DBDeleteRequest struct {
	TableName string        `json:"-"`
	Fields    []interface{} `json:"fields"`
	Data      []interface{} `json:"data"`
}

// DBDeleteResponse response for db delete data by id
type DBDeleteResponse struct {
	Status string `json:"status"`
}

func makeDBDeleteEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DBDeleteRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		rowData, relateRowData, err := sqlmapper.MakeRowData(req.Fields, req.Data)

		if err != nil {
			return nil, err
		}

		return DBDeleteResponse{"success"}, nil
	}
}
