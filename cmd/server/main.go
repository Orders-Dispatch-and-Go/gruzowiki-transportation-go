package main

import (
	"context"
	"gruzowiki/config"
	"gruzowiki/repositories"
	"gruzowiki/rest/handlers"
	"gruzowiki/rest/middlewares"
	"gruzowiki/services"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

type Server interface {
	Start()
}

type CarrierHandler interface {
	CreateCarrier(c echo.Context) error
	GetCarrier(c echo.Context) error
	UpdateCarrier(c echo.Context) error
	DeleteCarrier(c echo.Context) error
}

type CarHandler interface {
	CreateCar(c echo.Context) error
	GetCar(c echo.Context) error
	UpdateCar(c echo.Context) error
	DeleteCar(c echo.Context) error
	ListCarsByOwner(c echo.Context) error
}

type ServerImpl struct {
	Address        string
	CarrierHandler CarrierHandler
	CarHandler     CarHandler
}

func NewServer(address string, carrierHandler CarrierHandler, carHandler CarHandler) *ServerImpl {
	return &ServerImpl{
		Address:        address,
		CarrierHandler: carrierHandler,
		CarHandler:     carHandler,
	}
}

func (s *ServerImpl) Start() {
	e := echo.New()
	e.Use(middlewares.HandleError)

	carrier := e.Group("/carriers")
	carrier.POST("", s.CarrierHandler.CreateCarrier)
	carrier.GET("/:id", s.CarrierHandler.GetCarrier)
	carrier.PATCH("/:id", s.CarrierHandler.UpdateCarrier)
	carrier.DELETE("/:id", s.CarrierHandler.DeleteCarrier)

	cars := e.Group("/cars")
	cars.POST("", s.CarHandler.CreateCar)
	cars.GET("/:id", s.CarHandler.GetCar)
	cars.PATCH("/:id", s.CarHandler.UpdateCar)
	cars.DELETE("/:id", s.CarHandler.DeleteCar)
	cars.GET("/owners/:ownerId", s.CarHandler.ListCarsByOwner)

	e.Logger.Fatal(e.Start(s.Address))
}

func main() {
	cfg := config.MustLoad()

	logger := mustMakeLogger(cfg.LogLevel)
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

	carRepo := repositories.NewCarRepo(conn)
	carService := services.NewCarService(carRepo, carrierRepo)
	carHandler := handlers.NewCarHandler(carService)

	server := NewServer(cfg.Address, carrierHandler, carHandler)
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
