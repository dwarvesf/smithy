package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"

	"github.com/dwarvesf/smithy/backend"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	pgConfig "github.com/dwarvesf/smithy/backend/config/database/pg"
	"github.com/dwarvesf/smithy/backend/endpoints"
	serviceHttp "github.com/dwarvesf/smithy/backend/http"
	"github.com/dwarvesf/smithy/backend/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load .env by errors: %v", err))
	}
	var (
		httpAddr       = ":" + os.Getenv("PORT")
		configFilePath = os.Getenv("CONFIG_FILE_PATH")
	)

	cfg, err := backend.NewConfig(backendConfig.ReadYAML(configFilePath))
	if err != nil {
		panic(err)
	}

	ok, err := backend.SyncPersistent(cfg)
	if err != nil {
		panic(err)
	}

	if !ok {
		if err = cfg.UpdateConfigFromAgent(); err != nil {
			panic(err)
		}
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	pg, closeDB, err := pgConfig.NewPG()
	defer closeDB()
	if err != nil {
		panic(fmt.Sprintf("fail to connect to dashboard database by error %v", err))
	}

	if err = pgConfig.SeedCreateTable(pg); err != nil {
		panic(fmt.Sprintf("fail to migrate table by error %v", err))
	}

	s, err := service.NewService(cfg, pg)
	if err != nil {
		panic(err)
	}

	var h http.Handler
	{
		h = serviceHttp.NewHTTPHandler(
			endpoints.MakeServerEndpoints(s),
			logger,
			os.Getenv("ENV") == "local" || os.Getenv("ENV") == "development",
			cfg,
			s,
		)
	}
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errs <- http.ListenAndServe(httpAddr, h)
	}()

	_ = logger.Log("errors:", <-errs)
}
