package handler

import (
	"fmt"
	"net/http"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	handlerCommon "github.com/dwarvesf/smithy/common/handler"
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
