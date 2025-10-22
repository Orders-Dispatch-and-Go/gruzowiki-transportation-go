package handlers

import (
	"context"
	"gruzowiki/rest/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CarrierService interface {
	GetCarrier(context.Context, int32) (*models.GetCarrierResponse, error)
}

type Carrier struct {
	service CarrierService
}

func NewCarrierHandler(service CarrierService) *Carrier {
	return &Carrier{
		service: service,
	}
}

func (carriers *Carrier) GetCarrier(c echo.Context) error {
	idStr := c.Param("id")
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return err
	}

	carrier, err := carriers.service.GetCarrier(ctx, int32(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, carrier)
}
