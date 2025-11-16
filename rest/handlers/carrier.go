package handlers

import (
	"context"
	"fmt"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CarrierService interface {
	CreateCarrier(ctx context.Context, id int32, driverCategory string) (*models.CreateCarrierResponse, error)
	GetCarrier(context.Context, int32) (*models.GetCarrierResponse, error)
	UpdateCarrier(ctx context.Context, id int32, driverCategory string) (*models.UpdateCarrierResponse, error)
	DeleteCarrier(ctx context.Context, id int32) error
}

type CarrierHandler struct {
	service CarrierService
}

func NewCarrierHandler(service CarrierService) *CarrierHandler {
	return &CarrierHandler{service: service}
}

func (h *CarrierHandler) CreateCarrier(c echo.Context) error {
	var req models.CreateCarrierRequest
	if err := c.Bind(&req); err != nil {
		return terror.NewValidationError(err.Error())
	}
	fmt.Println("REQ:", req)

	resp, err := h.service.CreateCarrier(c.Request().Context(), req.ID, req.DriverCategory)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *CarrierHandler) GetCarrier(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return terror.NewValidationError("invalid id")
	}

	resp, err := h.service.GetCarrier(c.Request().Context(), int32(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *CarrierHandler) UpdateCarrier(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return terror.NewValidationError("invalid id")
	}

	var req models.UpdateCarrierRequest
	if err := c.Bind(&req); err != nil {
		return terror.NewValidationError(err.Error())
	}

	resp, err := h.service.UpdateCarrier(c.Request().Context(), int32(id), req.DriverCategory)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *CarrierHandler) DeleteCarrier(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return terror.NewValidationError("invalid id")
	}

	if err := h.service.DeleteCarrier(c.Request().Context(), int32(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
