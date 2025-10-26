package rest

import (
	"github.com/labstack/echo/v4"
	. "gruzowiki/rest/middlewares"
)

var (
	USER_ROLE    = "user"
	CARRIER_ROLE = "carrier"
	SENDER_ROLE  = "sender"
	ALL_ROLES    = []string{USER_ROLE, CARRIER_ROLE, SENDER_ROLE}
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

	e.Use(HandleError)

	carriers := e.Group("/carriers")
	carriers.GET("/:id", s.Carriers.GetCarrier, AllowedeRoles(ALL_ROLES...))

	startServer(e, s.Address)
}
