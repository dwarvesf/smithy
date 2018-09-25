package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
)

// DBUpdateRequest request for db Update data
type DBUpdateRequest struct {
	TableName    string        `json:"-"`
	DatabaseName string        `json:"-"`
	Fields       []interface{} `json:"fields"`
	Data         []interface{} `json:"data"`
}

// DBUpdateResponse response for db Update data
type DBUpdateResponse struct {
	Status string            `json:"status"`
	Data   sqlmapper.RowData `json:"data"`
}

func makeDBUpdateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DBUpdateRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		var err error
		rowData, err := sqlmapper.MakeRowData(req.Fields, req.Data)
		if err != nil {
			return nil, err
		}

		data, err := s.Update(req.DatabaseName, req.TableName, rowData)
		if err != nil {
			return nil, err
		}

		return DBUpdateResponse{"success", data}, nil
	}
}
