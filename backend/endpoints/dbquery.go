package endpoints

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend"
	"github.com/dwarvesf/smithy/backend/service"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

// DBQueryRequest request for db query
type DBQueryRequest struct {
	Method    string   `json:"method" schema:"method,required"`
	TableName string   `json:"table_name" schema:"table_name,required"`
	Cols      []string `json:"columns" schema:"columns,required"`
	Columns   []database.Column
	QueryData string `json:"query_data" schema:"query_data"`
	Offset    int    `json:"offset" schema:"offset" default:"0"`
	Limit     int    `json:"limit" schema:"limit" default:"-1"`
}

// DBQueryResponse response for db query
type DBQueryResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// UpdateColumnsByCols .
func (r *DBQueryRequest) UpdateColumnsByCols() error {
	res := []database.Column{}
	for _, col := range r.Cols {
		tmp := strings.Split(col, ",")
		if len(tmp) != 2 {
			return errors.New("wrong format of a column need at least 2 element")
		}

		name, colType := tmp[0], tmp[1]
		res = append(res, database.Column{Name: name, Type: colType})
	}

	r.Columns = res

	return nil
}

func (r *DBQueryRequest) getResourceID() (int, error) {
	return strconv.Atoi(r.QueryData)
}

func (r *DBQueryRequest) getColumnAndValue() (columnName string, value string, err error) {
	tmp := strings.Split(r.QueryData, ",")
	if len(tmp) != 2 {
		err = errors.New("query_data is wrong format")
		return
	}
	columnName = tmp[0]
	value = tmp[1]

	return
}

func makeDBQueryEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DBQueryRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}
		sqlmp, err := backend.NewSQLMapper(s.SyncConfig(), req.TableName, req.Columns)
		if err != nil {
			return nil, err
		}

		q := sqlmapper.Query{
			SourceTable: req.TableName,
			Fields:      req.Columns,
			Offset:      req.Offset,
			Limit:       req.Limit,
		}

		var data interface{}
		switch req.Method {
		case "FindByID":
			var id int
			if id, err = req.getResourceID(); err != nil {
				return nil, err
			}
			q.Filter.Value = strconv.Itoa(id)
			data, err = sqlmp.FindByID(q)
		case "FindAll":
			data, err = sqlmp.FindAll(q)
		case "FindByColumnName":
			var columnName, value string
			if columnName, value, err = req.getColumnAndValue(); err != nil {
				return nil, err
			}
			q.Filter = sqlmapper.Filter{ColName: columnName, Value: value}
			data, err = sqlmp.FindByColumnName(q)
		default:
			return nil, errors.New("unknown query method")
		}
		if err != nil {
			return nil, err
		}

		return DBQueryResponse{"success", data}, nil
	}
}
