package handlers

import (
	"context"
	"gruzowiki/repositories"
	"gruzowiki/rest/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CarrierService interface {
	GetCarrier(context.Context, string) (repositories.Carrier, error)
}

type Carrier struct {
	service CarrierService
}

func NewCarrier(service CarrierService) *Carrier {
	return &Carrier{
		service: service,
	}
}

func (carriers *Carrier) GetCarrier(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	carrier, err := carriers.service.GetCarrier(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.CarrierResponse{
		Id:             carrier.Id,
		DriverCategory: carrier.DriverCategory,
	})
}