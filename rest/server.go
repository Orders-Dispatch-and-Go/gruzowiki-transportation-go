package rest

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
	Address  string
	Carriers CarrierHandler
}

func NewServer(address string, carriers CarrierHandler) Server {
	return &ServerImpl{
		Address:  address,
		Carriers: carriers,
	}
}

func startServer(e *echo.Echo, address string) {
	e.Logger.Fatal(e.Start(address))
}

func (s *ServerImpl) Start() {
	e := echo.New()
	e.Use()

	e.HTTPErrorHandler = middlewares.ErrorHandler

	ping := e.Group("/carriers")
	ping.GET("/:id", s.Carriers.GetCarrier)

	startServer(e, s.Address)
}
