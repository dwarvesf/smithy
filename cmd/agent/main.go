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

	"github.com/dwarvesf/smithy/agent"
	"github.com/dwarvesf/smithy/agent/automigrate"
	agentConfig "github.com/dwarvesf/smithy/agent/config"
	agentHandler "github.com/dwarvesf/smithy/agent/handler"
)

var (
	httpAddr = ":" + os.Getenv("PORT")
)

func main() {
	cfg, err := agent.NewConfig(agentConfig.NewYAMLConfigReader("example_agent_config.yaml"))
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	err = automigrate.AutoMigrate(cfg)
	if err != nil {
		panic(err)
	}

	r.Get("/agent", agentHandler.Expose(cfg))

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
