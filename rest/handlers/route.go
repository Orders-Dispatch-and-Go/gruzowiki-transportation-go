package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"net/http"
)

type (
	RouteHandler struct {
		service RouteService
	}

	RouteService interface {
		GetRouteForCargoRequest(ctx context.Context, cargoRequestId uuid.UUID) (*models.GetTripRouteResponse, error)
	}
)

func NewRouteHandler(service RouteService) *RouteHandler {
	return &RouteHandler{service: service}
}

func (h *RouteHandler) GetRouteForCargoRequest(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("uuid")
	cargoRequestID, err := uuid.Parse(idParam)
	if err != nil {
		return terror.NewValidationError("invalid UUID", "parsing path parameter 'uuid'")
	}

	response, err := h.service.GetRouteForCargoRequest(ctx, cargoRequestID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
