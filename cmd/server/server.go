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

	ping := e.Group("")
	ping.GET("/carriers/:id", s.CarrierHandler.GetCarrier)
	ping.POST("/cargo_request", s.CargoRequestHandler.CreateCargoCargoRequest)

	startServer(e, s.Address)
}