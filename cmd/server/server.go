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

type ServerImpl struct {
	Address        string
	CarrierHandler CarrierHandler
}

func NewServer(address string, carrierHandler CarrierHandler) Server {
	return &ServerImpl{
		Address:        address,
		CarrierHandler: carrierHandler,
	}
}

func startServer(e *echo.Echo, address string) {
	e.Logger.Fatal(e.Start(address))
}

func (s *ServerImpl) Start() {
	e := echo.New()

	e.Use(middlewares.HandleError)

	ping := e.Group("/carriers")
	ping.GET("/:id", s.CarrierHandler.GetCarrier)

	startServer(e, s.Address)
}
