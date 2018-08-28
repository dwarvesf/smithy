package endpoints

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
)

// DBUpdateRequest request for db Update data
type DBUpdateRequest struct {
	TableName  string            `json:"table_name"`
	Data       sqlmapper.RowData `json:"data"`
	PrimaryKey string            `json:"primary_key" schema:"primary_key"`
}

// DBUpdateResponse response for db Update data
type DBUpdateResponse struct {
	Status string            `json:"status"`
	Data   sqlmapper.RowData `json:"data"`
}

func (r *DBUpdateRequest) getResourceID() (int, error) {
	return strconv.Atoi(r.PrimaryKey)
}

func makeDBUpdateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DBUpdateRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		var (
			id  int
			err error
		)
		if id, err = req.getResourceID(); err != nil {
			return nil, err
		}
		data, err := s.Update(req.TableName, req.Data, id)

		if err != nil {
			return nil, err
		}

		return DBQueryResponse{"success", data}, nil
	}
}
