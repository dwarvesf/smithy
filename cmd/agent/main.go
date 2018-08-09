package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwarvesf/smithy/config"
	agentConfig "github.com/dwarvesf/smithy/config/agent"
	"github.com/dwarvesf/smithy/config/agent/automigrate"
	"github.com/dwarvesf/smithy/handler/agent"
	"github.com/go-chi/chi"
)

var (
	httpAddr = ":" + os.Getenv("PORT")
)

func main() {
	cfg, err := config.NewAgentConfig(agentConfig.NewYAMLConfigReader("example_agent_config.yaml"))
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	err = automigrate.AutoMigrate(cfg)
	if err != nil {
		panic(err)
	}

	r.Get("/agent", agent.NewAgentHandler(cfg))

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
