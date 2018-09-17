package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	httpAddr = ":" + os.Getenv("PORT")
)

func main() {
	fs := http.FileServer(http.Dir("./swaggerui"))
	http.Handle("/", fs)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errs <- http.ListenAndServe(httpAddr, nil)
	}()

	log.Println("errors:", <-errs)
}
