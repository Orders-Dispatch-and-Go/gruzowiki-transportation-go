package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/gofiber/fiber/v2"
	"gruzowiki-transportation/internal/auth"
	"gruzowiki-transportation/internal/config"
	"gruzowiki-transportation/internal/middlewares"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type application struct {
	logger *log.Logger
}

func main() {
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := config.Load("")
	if err != nil {
		logger.Fatal("failed to load application configuration: %s", err)
		return
	}

	jwt := auth.NewJwt(cfg.Jwt)

	authMiddleware := []*middlewares.Middlewares{middlewares.New(jwt)}

	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(fiber.App{}, authMiddleware, logger, cfg),
	}

	go routing.GracefulShutdown(hs, 10*time.Second, logger.Printf)
	logger.Printf("server %s is running at %s", version, address)
	if err := hs.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
		return
	}
}

func buildHandler(
	app fiber.App,
	middlewares []*middlewares.Middlewares,
	logger *log.Logger,
	cfg config.Config,
) http.Handler {
	//app.Use(middlewares[0].Authenticate)
	//healthcheck.RegisterHandlers(router, version)

	app.Post("/healthcheck", middlewares[0].Authenticate)

	return nil
}
