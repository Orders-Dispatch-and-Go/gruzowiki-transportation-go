package rest

import (
	"github.com/labstack/echo/v4"
)

const (
	serverPrefix = "/gruzowiki/"
)

type Server interface {
	Start()
}

type СarriersHandler interface {
	GetCarrier(c echo.Context) error
}

type ServerImpl struct {
	Address  string
	Сarriers СarriersHandler
}

func NewServer(address string, carriers СarriersHandler) Server {
	return &ServerImpl{
		Address:  address,
		Сarriers: carriers,
	}
}

func startServer(e *echo.Echo, address string) {
	e.Logger.Fatal(e.Start(address))
}

func (s *ServerImpl) Start() {
	e := echo.New()
	e.Use()

	gruzowiki := e.Group(serverPrefix)

	ping := gruzowiki.Group("/Сarriers")
	ping.GET("/:id", s.Сarriers.GetCarrier)

	startServer(e, s.Address)
}
