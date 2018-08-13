package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/schema"

	"github.com/dwarvesf/smithy/backend"
	"github.com/dwarvesf/smithy/common/database"
	handlerCommon "github.com/dwarvesf/smithy/common/handler"
)

// QueryRequest request for query data
// TODO: verify query request query data for matching with each find
type QueryRequest struct {
	Method    string   `json:"method" schema:"method,required"`
	TableName string   `json:"table_name" schema:"table_name,required"`
	Cols      []string `json:"columns" schema:"columns,required"`
	Columns   []database.Column
	QueryData string `json:"query_data" schema:"query_data"`
}

func (r *QueryRequest) updateColumnsByCols() error {
	res := []database.Column{}
	for _, col := range r.Cols {
		tmp := strings.Split(col, ",")
		if len(tmp) != 2 {
			return errors.New("Wrong format of a column need at least 2 element")
		}

		name, colType := tmp[0], tmp[1]
		res = append(res, database.Column{Name: name, Type: colType})
	}

	r.Columns = res

	return nil
}

// Query query request
func (h *Handler) Query() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
			return
		}

		qr := new(QueryRequest)
		decoder := schema.NewDecoder()
		err = decoder.Decode(qr, r.Form)
		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
			return
		}

		qr.updateColumnsByCols()

		sqlmp, err := backend.NewSQLMapper(h.cfg.Config(), qr.TableName, qr.Columns)
		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
			return
		}

		var buf []byte
		switch qr.Method {
		case "FindByID":
			var id int
			if id, err = strconv.Atoi(qr.QueryData); err != nil {
				handlerCommon.EncodeJSONError(err, w)
				return
			}
			buf, err = sqlmp.FindByID(id)
		case "FindAll":
			buf, err = sqlmp.FindAll()
		default:
			handlerCommon.EncodeJSONError(errors.New("Unknown query method"), w)
			return
		}

		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
			return
		}

		handlerCommon.EncodeJSONSuccess(buf, w)
	}
}
