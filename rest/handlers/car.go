package handlers

import (
	"context"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CarService interface {
	CreateCar(ctx context.Context, req models.CreateCarRequest) (*models.CreateCarResponse, error)
	GetCar(ctx context.Context, id int32) (*models.GetCarResponse, error)
	UpdateCar(ctx context.Context, id int32, req models.UpdateCarRequest) (*models.UpdateCarResponse, error)
	DeleteCar(ctx context.Context, id int32) error
	ListCarsByOwner(ctx context.Context, ownerID int32) ([]*models.GetCarResponse, error)
}

type CarHandler struct {
	service CarService
}

func NewCarHandler(service CarService) *CarHandler {
	return &CarHandler{service: service}
}

func (h *CarHandler) CreateCar(c echo.Context) error {
	var req models.CreateCarRequest
	if err := c.Bind(&req); err != nil {
		return terror.NewValidationError(err.Error())
	}

	resp, err := h.service.CreateCar(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *CarHandler) GetCar(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return terror.NewValidationError("invalid id")
	}

	resp, err := h.service.GetCar(c.Request().Context(), int32(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *CarHandler) UpdateCar(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return terror.NewValidationError("invalid id")
	}

	var req models.UpdateCarRequest
	if err := c.Bind(&req); err != nil {
		return terror.NewValidationError(err.Error())
	}

	resp, err := h.service.UpdateCar(c.Request().Context(), int32(id), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *CarHandler) DeleteCar(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return terror.NewValidationError("invalid id")
	}

	if err := h.service.DeleteCar(c.Request().Context(), int32(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CarHandler) ListCarsByOwner(c echo.Context) error {
	ownerID, err := strconv.Atoi(c.Param("ownerId"))
	if err != nil {
		return terror.NewValidationError("invalid ownerId")
	}

	cars, err := h.service.ListCarsByOwner(c.Request().Context(), int32(ownerID))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cars)
}
