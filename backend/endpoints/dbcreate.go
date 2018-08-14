package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend"
	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

// DBCreateRequest request for db create data
type DBCreateRequest struct {
	TableName string            `json:"table_name"`
	Data      sqlmapper.RowData `json:"data"`
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

		sqlmp, err := backend.NewSQLMapper(s.Config.Config(), req.TableName, []database.Column{})
		if err != nil {
			return nil, err
		}

		data, err := sqlmp.Create(req.Data)
		if err != nil {
			return nil, err
		}

		return DBQueryResponse{"success", data}, nil
	}
}
