package handler

import (
	"fmt"
	"net/http"

	"github.com/dwarvesf/smithy/backend"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
	handlerCommon "github.com/dwarvesf/smithy/common/handler"
	"github.com/k0kubun/pp"
)

// Handler handler for dashboard
type Handler struct {
	cfg *backendConfig.Wrapper
}

// NewHandler new dashboard handler
func NewHandler(cfg *backendConfig.Config) *Handler {
	return &Handler{backendConfig.NewWrapper(cfg)}
}

// NewUpdateConfigFromAgent return handler for expose metadata, connection for dashboard
func (h *Handler) NewUpdateConfigFromAgent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := h.cfg.Config()
		err := cfg.UpdateConfigFromAgent()
		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
			return
		}

		fmt.Fprintln(w, `{"status": "success"}`)
	}
}

// NewCRUD .
// FIXME: REMOVE
// TODO: REMOVE and UPDATE FOR TMP ONLY
func (h *Handler) NewCRUD() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlmp, err := backend.NewSQLMapper(h.cfg.Config(), "users", []database.Column{
			{
				Name: "id",
				Type: "int",
			},
			{
				Name: "name",
				Type: "string",
			},
			{
				Name: "age",
				Type: "int",
			},
		})
		if err != nil {
			handlerCommon.EncodeJSONError(err, w)
			return
		}
		buf, err := sqlmp.Create(map[string]sqlmapper.ColData{
			"name": sqlmapper.ColData{Name: "name", Data: "hieu"},
			"age":  sqlmapper.ColData{Name: "name", Data: 26},
		})
		if err != nil {
			fmt.Fprintln(w, `{"status": "failure"}`)
			pp.Print(err)
			return
		}
		pp.Print(string(buf))

		fmt.Fprintln(w, `{"status": "success"}`)
	}
}
