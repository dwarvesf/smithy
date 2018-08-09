package dashboard

import (
	"fmt"
	"net/http"

	"github.com/dwarvesf/smithy/config/dashboard"
	"github.com/dwarvesf/smithy/handler/common"
)

// Handler handler for dashboard
type Handler struct {
	Config *dashboard.Config
}

// NewDashboardHandler new dashboard handler
func NewDashboardHandler(cfg *dashboard.Config) *Handler {
	return &Handler{cfg}
}

// NewUpdateConfigFromAgent return handler for expose metadata, connection for dashboard
func (h *Handler) NewUpdateConfigFromAgent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.Config.UpdateConfigFromAgent()
		if err != nil {
			common.EncodeJSONError(err, w)
			return
		}

		fmt.Fprintln(w, `{"status": "success"}`)
	}
}
