package main

import (
	"fmt"
	"gruzowiki/config"
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
	Delivered(c echo.Context) error
	GetRequestsForTrip(c echo.Context) error
}

type TripHandler interface {
	GetTripByCargoRequest(c echo.Context) error
	CreateTrip(c echo.Context) error
	Finish(c echo.Context) error
	Start(c echo.Context) error
	GetTripByIdAndCarrier(c echo.Context) error
}

type RoutesHandler interface {
	GetRouteForCargoRequest(c echo.Context) error
}

type ServerImpl struct {
	Address             string
	Config              config.Config
	CarrierHandler      CarrierHandler
	CargoRequestHandler CargoRequestHandler
	CarHandler          CarHandler
	RecipientHandler    RecipientHandler
	TripHandler         TripHandler
	ConsignerHandler    ConsignerHandler
	RoutesHandler       RoutesHandler
}

func NewServer(
	address string,
	config config.Config,
	carrierHandler CarrierHandler,
	cargoRequestHandler CargoRequestHandler,
	carHandler CarHandler,
	recipientHandler RecipientHandler,
	tripHandler TripHandler,
	consignerHandler ConsignerHandler,
	routesHandler RoutesHandler,
) Server {
	return &ServerImpl{
		Address:             address,
		Config:              config,
		CarrierHandler:      carrierHandler,
		CargoRequestHandler: cargoRequestHandler,
		CarHandler:          carHandler,
		RecipientHandler:    recipientHandler,
		TripHandler:         tripHandler,
		ConsignerHandler:    consignerHandler,
		RoutesHandler:       routesHandler,
	}
}

func (s *ServerImpl) startServer(e *echo.Echo, address string) {
	fmt.Printf("server configuration:")
	fmt.Printf("\nauth client: %s\nrust client: %s\n", s.Config.AuthClientUrl, s.Config.RustClientUrl)
	e.Logger.Fatal(e.Start(address))
}

func (s *ServerImpl) Start() {
	e := echo.New()
	//e.Use(middlewares.LoggingMiddleware)
	e.Use(middlewares.HandleError)

	postCarrier := e.Group("")
	postCarrier.POST("/carrier", s.CarrierHandler.CreateCarrier)

	carriers := e.Group("")
	carriers.Use(middlewares.AllowedRoles([]string{middlewares.ConsignerRole, middlewares.CarrierRole}...))
	carriers.GET("/carrier/:id", s.CarrierHandler.GetCarrier)
	//carriers.PUT("/carrier/:id", s.CarrierHandler.UpdateCarrier)
	//carriers.DELETE("/carrier/:id", s.CarrierHandler.DeleteCarrier)

	cargoRequest := e.Group("/cargo_request")
	cargoRequest.Use(middlewares.AllowedRoles([]string{middlewares.ConsignerRole, middlewares.CarrierRole}...))
	cargoRequest.POST("/search", s.CargoRequestHandler.GetCargoRequest)
	cargoRequest.POST("", s.CargoRequestHandler.CreateCargoCargoRequest)
	cargoRequest.POST("/:cargo_request/:cargoRequestId/trip/:tripId", s.CargoRequestHandler.MarkTrip)
	cargoRequest.PATCH("/:uuid/finish/code/:code", s.CargoRequestHandler.Delivered)
	cargoRequest.GET("/trip/:tripID", s.CargoRequestHandler.GetRequestsForTrip)

	cargo := e.Group("/cargo")
	cargo.Use(middlewares.AllowedRoles([]string{middlewares.ConsignerRole, middlewares.CarrierRole}...))
	cargo.GET("/types", s.CargoRequestHandler.GetCargoTypes)
	cargo.POST("", s.CargoRequestHandler.CreateCargo)

	cars := e.Group("/cars")
	cars.Use(middlewares.AllowedRoles([]string{middlewares.ConsignerRole, middlewares.CarrierRole}...))
	cars.POST("", s.CarHandler.CreateCar)
	cars.GET("/:id", s.CarHandler.GetCar)
	//cars.PUT("/:id", s.CarHandler.UpdateCar)
	//cars.DELETE("/:id", s.CarHandler.DeleteCar)
	cars.GET("/owner/:ownerId", s.CarHandler.ListCarsByOwner)

	recipients := e.Group("/recipients")
	recipients.Use(middlewares.AllowedRoles([]string{middlewares.ConsignerRole, middlewares.CarrierRole}...))
	recipients.POST("", s.RecipientHandler.CreateRecipient)
	recipients.GET("/:id", s.RecipientHandler.GetRecipient)
	recipients.GET("", s.RecipientHandler.ListRecipients)
	recipients.PUT("/:id", s.RecipientHandler.UpdateRecipient)
	recipients.DELETE("/:id", s.RecipientHandler.DeleteRecipient)

	consigners := e.Group("/consigner")
	consigners.POST("", s.ConsignerHandler.CreateConsigner)

	trips := e.Group("/trip")
	trips.GET("", s.TripHandler.GetTripByIdAndCarrier)
	trips.GET("/cargo_request/:id", s.TripHandler.GetTripByCargoRequest)
	trips.POST("", s.TripHandler.CreateTrip)
	trips.PATCH("/:id/finish/status/:status", s.TripHandler.Finish)
	trips.PATCH("/:id/start", s.TripHandler.Start)

	routes := e.Group("/routes")
	routes.Use(middlewares.AllowedRoles([]string{middlewares.ConsignerRole, middlewares.CarrierRole}...))
	routes.GET("/cargo_request/:uuid", s.RoutesHandler.GetRouteForCargoRequest)

	s.startServer(e, s.Address)
}
