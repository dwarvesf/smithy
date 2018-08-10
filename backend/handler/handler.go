package handler

import (
	"fmt"
	"net/http"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	pgMapper "github.com/dwarvesf/smithy/backend/sqlmapper/pg"
	"github.com/dwarvesf/smithy/common/database"
	handlerCommon "github.com/dwarvesf/smithy/common/handler"
	"github.com/k0kubun/pp"
)

// Handler handler for dashboard
type Handler struct {
	Config *backendConfig.Config
}

// NewHandler new dashboard handler
func NewHandler(cfg *backendConfig.Config) *Handler {
	return &Handler{cfg}
}

// NewUpdateConfigFromAgent return handler for expose metadata, connection for dashboard
func (h *Handler) NewUpdateConfigFromAgent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.Config.UpdateConfigFromAgent()
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

		buf, err := pgMapper.NewQuerier(h.Config.GetDB(), "users", []database.Column{
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
		}).Create(map[string]sqlmapper.ColData{
			"name": sqlmapper.ColData{Name: "name", Data: "hieu"},
			"age":  sqlmapper.ColData{Name: "name", Data: 26},
		})
		if err != nil {
			pp.Print(err)
		}
		pp.Print(string(buf))

		fmt.Fprintln(w, `{"status": "success"}`)
	}
}
