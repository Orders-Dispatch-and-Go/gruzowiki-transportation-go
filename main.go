package main

import (
	"gruzowiki/config"
	"gruzowiki/repositories"
	"gruzowiki/rest"
	"gruzowiki/rest/handlers"
	"gruzowiki/services"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	loger := mustMakeLogger(cfg.LogLevel)
	//TO DO прокинуть везде логер

	loger.Info("starting server")

	conn, err := repositories.NewConnect(cfg.Address)
	if err != nil {
		loger.Error(err.Error())
		return
	}

	carrierRepo := repositories.NewCarrier(conn)

	carrierService := services.NewCarrierService(carrierRepo)

	carrierHandler := handlers.NewCarrier(carrierService)

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