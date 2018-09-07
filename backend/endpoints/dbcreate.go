package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
)

// DBCreateRequest request for db create data
type DBCreateRequest struct {
	TableName string        `json:"-"`
	Fields    []interface{} `json:"fields"`
	Data      []interface{} `json:"data"`
}

// DBCreateResponse response for db create data
type DBCreateResponse struct {
	Status string            `json:"status"`
	Data   sqlmapper.RowData `json:"data"`
}

func makeDBCreateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DBCreateRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		rowData, err := sqlmapper.MakeRowData(req.Fields, req.Data)
		if err != nil {
			return nil, err
		}

		data, err := s.Create(req.TableName, rowData)
		if err != nil {
			return nil, err
		}

		return DBCreateResponse{"success", data}, nil
	}
}
