package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/dwarvesf/smithy/backend"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	backendHandler "github.com/dwarvesf/smithy/backend/handler"
)

var (
	httpAddr = ":" + os.Getenv("PORT")
)

func main() {
	cfg, err := backend.NewConfig(backendConfig.ReadYAML("example_dashboard_config.yaml"))
	if err != nil {
		panic(err)
	}

	h := backendHandler.NewHandler(cfg)

	r := chi.NewRouter()

	r.Get("/agent-sync", h.NewUpdateConfigFromAgent())
	// Query
	r.Get("/query", h.Query())
	r.Post("/query", h.Query()) // Post query for query parameter longer than 2048 character

	// Create
	r.Post("/create", h.Create())

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
