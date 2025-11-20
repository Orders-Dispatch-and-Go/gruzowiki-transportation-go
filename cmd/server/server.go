package main

import (
	"github.com/labstack/echo/v4"
	"gruzowiki/rest/middlewares"
)

type Server interface {
	Start()
}

type CarrierHandler interface {
	GetCarrier(c echo.Context) error
}

type CargoRequestHandler interface {
	CreateCargoCargoRequest(c echo.Context) error
	GetCargoTypes(c echo.Context) error
	CreateCargo(c echo.Context) error
}

type ServerImpl struct {
	Address             string
	CarrierHandler      CarrierHandler
	CargoRequestHandler CargoRequestHandler
}

func NewServer(address string, carrierHandler CarrierHandler, cargoRequestHandler CargoRequestHandler) Server {
	return &ServerImpl{
		Address:             address,
		CarrierHandler:      carrierHandler,
		CargoRequestHandler: cargoRequestHandler,
	}
}

func startServer(e *echo.Echo, address string) {
	e.Logger.Fatal(e.Start(address))
}

func (s *ServerImpl) Start() {
	e := echo.New()

	e.Use(middlewares.HandleError)

	carriers := e.Group("/carriers")
	carriers.GET("/:id", s.CarrierHandler.GetCarrier)

	cargoRequest := e.Group("/cargo_request")
	cargoRequest.POST("/cargo_request", s.CargoRequestHandler.CreateCargoCargoRequest)

	cargo := e.Group("/cargo")
	cargo.GET("/types", s.CargoRequestHandler.GetCargoTypes)
	cargo.POST("", s.CargoRequestHandler.CreateCargo)

	startServer(e, s.Address)
}
