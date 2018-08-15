package endpoints

import (
	"context"
	"errors"
	"strconv"

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
	QueryData string            `json:"primary_id" schema:"primary_id"`
}

// DBUpdateResponse response for db Update data
type DBUpdateResponse struct {
	Status string            `json:"status"`
	Data   sqlmapper.RowData `json:"data"`
}

func (r *DBUpdateRequest) getResourceID() (int, error) {
	return strconv.Atoi(r.QueryData)
}

func makeDBUpdateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DBUpdateRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		sqlmp, err := backend.NewSQLMapper(s.Config.Config(), req.TableName, []database.Column{})
		if err != nil {
			return nil, err
		}

		var id int
		if id, err = req.getResourceID(); err != nil {
			return nil, err
		}
		data, err := sqlmp.Update(req.Data, id)

		if err != nil {
			return nil, err
		}

		return DBQueryResponse{"success", data}, nil
	}
}
