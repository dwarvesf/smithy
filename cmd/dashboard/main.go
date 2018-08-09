package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwarvesf/smithy/config"
	dashboardConfig "github.com/dwarvesf/smithy/config/dashboard"
	dashboardHandler "github.com/dwarvesf/smithy/handler/dashboard"
	"github.com/go-chi/chi"
)

var (
	httpAddr = ":" + os.Getenv("PORT")
)

func main() {
	cfg, err := config.NewDashboardConfig(dashboardConfig.NewYAMLConfigReader("example_dashboard_config.yaml"))
	if err != nil {
		panic(err)
	}

	h := dashboardHandler.NewDashboardHandler(cfg)

	r := chi.NewRouter()

	r.Get("/agent-sync", h.NewUpdateConfigFromAgent())

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errs <- http.ListenAndServe(httpAddr, r)
	}()

	log.Println("errors:", <-errs)
}
