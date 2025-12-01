package main

import (
	"context"
	"gruzowiki/client"
	"gruzowiki/config"
	"gruzowiki/repositories"
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

	client := client.NewRoutesClient(cfg.ClinetUrl)

	carrierRepo := repositories.NewCarrierRepo(conn)
	carrierService := services.NewCarrierService(carrierRepo)
	carrierHandler := handlers.NewCarrierHandler(carrierService)

	stationRepo := repositories.NewStationRepo(conn)
	stationService := services.NewStationService(stationRepo)

	consignerRepo := repositories.NewConsignerRepo(conn)
	consignerService := services.NewConsignerService(consignerRepo)
	consignerHandler := handlers.NewConsignerHandler(consignerService)

	cargoRequestRepo := repositories.NewCargoRequestRepo(conn)
	cargoRequestService := services.NewCargoRequestService(cargoRequestRepo, stationService)
	cargoRequestHandler := handlers.NewCargoRequestController(cargoRequestService)

	carRepo := repositories.NewCarRepo(conn)
	carService := services.NewCarService(carRepo, carrierRepo)
	carHandler := handlers.NewCarHandler(carService)

	recipientRepo := repositories.NewRecipientRepo(conn)
	recipientService := services.NewRecipientService(recipientRepo)
	recipientHandler := handlers.NewRecipientHandler(recipientService)

	tripRepo := repositories.NewTripRepo(conn)
	tripService := services.NewTripService(tripRepo, stationRepo, client)
	tripHandler := handlers.NewTripHandler(tripService)

	server := NewServer(cfg.Address, carrierHandler, cargoRequestHandler, carHandler, recipientHandler, tripHandler, consignerHandler)
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