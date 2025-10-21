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

type СarriersHandler interface{
	GetCarrier(c echo.Context) error
}

type ServerImpl struct {
	Port string
	Сarriers СarriersHandler
}

func NewServer(port string, carriers СarriersHandler) Server {
	return &ServerImpl{
		Port: port,
		Сarriers: carriers,
	}
}

func startServer(e *echo.Echo, port string) {
	e.Logger.Fatal(e.Start("127.0.0.1:" + port))
}

func (s *ServerImpl) Start() {
	e := echo.New()
	e.Use()

	gruzowiki:= e.Group(serverPrefix)

	ping := gruzowiki.Group("/Сarriers")
	ping.GET("/:id", s.Сarriers.GetCarrier)

	startServer(e, s.Port)
}