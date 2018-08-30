package endpoints

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
)

// DBDeleteRequest request for db delete data by id
type DBDeleteRequest struct {
	TableName  string `json:"-"`
	PrimaryKey string `json:"primary_key" schema:"primary_key"`
}

// DBDeleteResponse response for db delete data by id
type DBDeleteResponse struct {
	Status string `json:"status"`
}

func (r *DBDeleteRequest) getResourceID() (int, error) {
	return strconv.Atoi(r.PrimaryKey)
}

func makeDBDeleteEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DBDeleteRequest)
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

		if err := s.Delete(req.TableName, id); err != nil {
			return nil, err
		}

		return DBDeleteResponse{"success"}, nil
	}
}
