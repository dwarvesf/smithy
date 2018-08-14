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

// DBUpdateRequest request for db Update data
type DBUpdateRequest struct {
	TableName string            `json:"table_name"`
	Data      sqlmapper.RowData `json:"data"`
	QueryData string            `json:"query_data" schema:"query_data"`
}

// DBUpdateResponse response for db Update data
type DBUpdateResponse struct {
	Status string            `json:"status"`
	Data   sqlmapper.RowData `json:"data"`
}

func (r *DBUpdateRequest) getResourceID() string {
	return r.QueryData
}

func makeDBUpdateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DBUpdateRequest)

		id := req.getResourceID()

		if id == "" {
			return nil, errors.New("id is not exist")
		}
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		sqlmp, err := backend.NewSQLMapper(s.Config.Config(), req.TableName, []database.Column{})
		if err != nil {
			return nil, err
		}

		data, err := sqlmp.Update(req.Data, id)
		if err != nil {
			return nil, err
		}

		return DBQueryResponse{"success", data}, nil
	}
}
