package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
)

type deleteFilter struct {
	Fields []interface{} `json:"fields"`
	Data   []interface{} `json:"data"`
}

// DBDeleteRequest request for db delete data by id
type DBDeleteRequest struct {
	TableName    string       `json:"-"`
	DatabaseName string       `json:"-"`
	Filter       deleteFilter `json:"filter"`
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

		if err := s.Delete(req.DatabaseName, req.TableName, req.Filter.Fields, req.Filter.Data); err != nil {
			return nil, err
		}

		return DBDeleteResponse{"success"}, nil
	}
}
