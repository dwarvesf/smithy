package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/dwarvesf/smithy/backend"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/endpoints"
	serviceHttp "github.com/dwarvesf/smithy/backend/http"
	"github.com/dwarvesf/smithy/backend/service"
)

var (
	httpAddr = ":" + os.Getenv("PORT")
)

func main() {
	cfg, err := backend.NewConfig(backendConfig.ReadYAML("example_dashboard_config.yaml"))
	if err != nil {
		panic(err)
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	s := service.NewService(cfg)

<<<<<<< HEAD
	var h http.Handler
	{
		h = serviceHttp.NewHTTPHandler(
			endpoints.MakeServerEndpoints(s),
			logger,
			os.Getenv("ENV") == "local",
		)
	}
=======
	r.Get("/agent-sync", h.NewUpdateConfigFromAgent())
	//r.Get("/crud", h.NewCRUD()) // TODO: REMOVE THIS
	r.Put("/update/{id}", h.Update()) // TODO: REMOVE THIS
>>>>>>> implements update

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errs <- http.ListenAndServe(httpAddr, h)
	}()

	logger.Log("errors:", <-errs)
}
