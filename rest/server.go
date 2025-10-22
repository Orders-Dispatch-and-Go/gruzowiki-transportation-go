package rest

import (
	"github.com/labstack/echo/v4"
)

type Server interface {
	Start()
}

type СarriersHandler interface {
	GetCarrier(c echo.Context) error
}

type ServerImpl struct {
	Address  string
	Carriers СarriersHandler
}

func NewServer(address string, carriers СarriersHandler) Server {
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

	ping := e.Group("/carriers")
	ping.GET("/:id", s.Carriers.GetCarrier)

	startServer(e, s.Address)
}
