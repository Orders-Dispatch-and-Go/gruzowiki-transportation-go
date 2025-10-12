package main

import (
	"errors"
	"flag"
	"fmt"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"gruzowiki-transportation/internal/config"
	"gruzowiki-transportation/internal/healthcheck"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type application struct {
	logger *log.Logger
}

var flagConfig = flag.String("config", "./config/local.yaml", "path to the config file")

func main() {
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Fatal("failed to load application configuration: %s", err)
		return
	}

	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, cfg),
	}

	go routing.GracefulShutdown(hs, 10*time.Second, logger.Printf)
	logger.Printf("server %s is running at %s", version, address)
	if err := hs.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
		return
	}
}

func buildHandler(logger *log.Logger, cfg *config.Config) http.Handler {
	router := routing.New()

	healthcheck.RegisterHandlers(router, version)

	return router
}
