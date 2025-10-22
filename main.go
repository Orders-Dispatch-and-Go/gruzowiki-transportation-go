package main

import (
	"context"
	"gruzowiki/config"
	"gruzowiki/repositories"
	"gruzowiki/rest"
	"gruzowiki/rest/handlers"
	"gruzowiki/services"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	logger := mustMakeLogger(cfg.LogLevel)
	//TO DO прокинуть везде логер

	logger.Info("starting server")

	ctx, stopCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopCtx()

	conn, err := repositories.NewConnect(ctx, cfg.Dsn)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	carrierRepo := repositories.NewCarrierRepo(conn)
	carrierService := services.NewCarrierService(carrierRepo)
	carrierHandler := handlers.NewCarrierHandler(carrierService)

	server := rest.NewServer(cfg.Address, carrierHandler)
	server.Start()
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		panic("unknown log level: " + logLevel)
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
