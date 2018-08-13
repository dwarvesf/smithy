package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dwarvesf/smithy/backend"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
	handlerCommon "github.com/dwarvesf/smithy/common/handler"
)

// CreateRequest request for create data
type CreateRequest struct {
	TableName string            `json:"table_name"`
	Data      sqlmapper.RowData `json:"data"`
}

// Create request for create a data
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
		}
		defer r.Body.Close()

		var cr CreateRequest
		json.Unmarshal(buf, &cr)

		sqlmp, err := backend.NewSQLMapper(h.cfg.Config(), cr.TableName, []database.Column{})
		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
			return
		}

		buf, err = sqlmp.Create(cr.Data)
		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
			return
		}

		handlerCommon.EncodeJSONSuccess(buf, w)
	}
}
