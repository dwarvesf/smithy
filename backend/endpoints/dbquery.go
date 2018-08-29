package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

// DBQueryRequest request for db query
type DBQueryRequest struct {
	sqlmapper.Query
}

// DBQueryResponse response for db query
type DBQueryResponse struct {
	Status  string            `json:"status"`
	Columns []string          `json:"columns,omitempty"`
	Rows    []interface{}     `json:"rows,omitempty"`
	Cols    []database.Column `json:"cols,omitempty"`
}

// // UpdateColumnsByCols .
// func (r *DBQueryRequest) UpdateColumnsByCols() error {
// 	res := []database.Column{}
// 	for _, col := range r.Cols {
// 		tmp := strings.Split(col, ",")
// 		if len(tmp) != 2 {
// 			return errors.New("wrong format of a column need at least 2 element")
// 		}

// 		name, colType := tmp[0], tmp[1]
// 		res = append(res, database.Column{Name: name, Type: colType})
// 	}

// 	r.Columns = res

// 	return nil
// }

func makeDBQueryEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DBQueryRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		columns, data, err := s.Query(req.Query)
		if err != nil {
			return nil, err
		}

		return DBQueryResponse{
			Status:  "success",
			Columns: columns,
			Rows:    data,
		}, nil
	}
}
