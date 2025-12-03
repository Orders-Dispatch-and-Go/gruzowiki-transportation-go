package main

import (
	"gruzowiki/rest/middlewares"

	"github.com/labstack/echo/v4"
)

type Server interface {
	Start()
}

type ConsignerHandler interface {
	CreateConsigner(c echo.Context) error
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

type RecipientHandler interface {
	CreateRecipient(c echo.Context) error
	GetRecipient(c echo.Context) error
	ListRecipients(c echo.Context) error
	UpdateRecipient(c echo.Context) error
	DeleteRecipient(c echo.Context) error
}

type CargoRequestHandler interface {
	GetCargoRequest(c echo.Context) error
	CreateCargoCargoRequest(c echo.Context) error
	GetCargoTypes(c echo.Context) error
	CreateCargo(c echo.Context) error
	MarkTrip(c echo.Context) error
}

type TripHandler interface {
	GetTripByCargoRequest(c echo.Context) error
	CreateTrip(c echo.Context) error
}

type ServerImpl struct {
	Address             string
	CarrierHandler      CarrierHandler
	CargoRequestHandler CargoRequestHandler
	CarHandler          CarHandler
	RecipientHandler    RecipientHandler
	TripHandler         TripHandler
	ConsignerHandler    ConsignerHandler
}

func NewServer(address string, carrierHandler CarrierHandler, cargoRequestHandler CargoRequestHandler, carHandler CarHandler, recipientHandler RecipientHandler,
	tripHandler TripHandler, consignerHandler ConsignerHandler) Server {
	return &ServerImpl{
		Address:             address,
		CarrierHandler:      carrierHandler,
		CargoRequestHandler: cargoRequestHandler,
		CarHandler:          carHandler,
		RecipientHandler:    recipientHandler,
		TripHandler:         tripHandler,
		ConsignerHandler:    consignerHandler,
	}
}

func startServer(e *echo.Echo, address string) {
	e.Logger.Fatal(e.Start(address))
}

func (s *ServerImpl) Start() {
	e := echo.New()

	e.Use(middlewares.HandleError)

	e.Use(middlewares.AllowedRoles([]string{middlewares.ConsignerRole, middlewares.CarrierRole}...))

	carriers := e.Group("/carrier")
	carriers.GET("/:id", s.CarrierHandler.GetCarrier)
	carriers.POST("", s.CarrierHandler.CreateCarrier) //здесь не должно быть проверки на jwt токен
	carriers.PUT("/:id", s.CarrierHandler.UpdateCarrier)
	carriers.DELETE("/:id", s.CarrierHandler.DeleteCarrier)

	cargoRequest := e.Group("/cargo_request")
	cargoRequest.POST("/search", s.CargoRequestHandler.GetCargoRequest)
	cargoRequest.POST("", s.CargoRequestHandler.CreateCargoCargoRequest)
	cargoRequest.POST("/:cargo_requestId/trip/tripId", s.CargoRequestHandler.MarkTrip)

	cargo := e.Group("/cargo")
	cargo.GET("/types", s.CargoRequestHandler.GetCargoTypes)
	cargo.POST("", s.CargoRequestHandler.CreateCargo)

	cars := e.Group("/cars")
	cars.POST("", s.CarHandler.CreateCar)
	cars.GET("/:id", s.CarHandler.GetCar)
	cars.PUT("/:id", s.CarHandler.UpdateCar)
	cars.DELETE("/:id", s.CarHandler.DeleteCar)
	cars.GET("/owner/:ownerId", s.CarHandler.ListCarsByOwner)

	recipients := e.Group("/recipients")
	recipients.POST("", s.RecipientHandler.CreateRecipient)
	recipients.GET("/:id", s.RecipientHandler.GetRecipient)
	recipients.GET("", s.RecipientHandler.ListRecipients)
	recipients.PUT("/:id", s.RecipientHandler.UpdateRecipient)
	recipients.DELETE("/:id", s.RecipientHandler.DeleteRecipient)

	consigners := e.Group("/consigner")
	consigners.POST("", s.ConsignerHandler.CreateConsigner) //здесь не должно быть проверки на jwt токен

	trips := e.Group("/trip")
	trips.GET("/cargo_request/:id", s.TripHandler.GetTripByCargoRequest)
	trips.POST("", s.TripHandler.CreateTrip)

	startServer(e, s.Address)
}
