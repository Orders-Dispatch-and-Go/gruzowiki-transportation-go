package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gruzowiki/rest/middlewares"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"net/http"
	"strconv"
)

type (
	RouteHandler struct {
		service RouteService
	}

	RouteService interface {
		GetRouteForCargoRequest(ctx context.Context, cargoRequestId uuid.UUID, withPoints bool) (*models.GetTripRouteResponse, error)
		GetRouteForTrip(ctx context.Context, tripId uuid.UUID, withPoints bool) (*models.GetTripRouteResponse, error)
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

	withPointsParam := c.QueryParam("withPoints")
	if withPointsParam == "" {
		withPointsParam = "false"
	}

	withPoints, err := strconv.ParseBool(withPointsParam)
	if err != nil {
		return terror.NewValidationError("bad validate", "withPoints")
	}

	userId := c.Get(middlewares.UserIdCtxClaim)
	ctx = context.WithValue(ctx, middlewares.UserIdCtxClaim, userId)

	response, err := h.service.GetRouteForCargoRequest(ctx, cargoRequestID, withPoints)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *RouteHandler) GetRouteForTrip(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("uuid")
	tripId, err := uuid.Parse(idParam)
	if err != nil {
		return terror.NewValidationError("invalid UUID", "parsing path parameter 'uuid'")
	}

	withPointsParam := c.QueryParam("withPoints")
	if withPointsParam == "" {
		withPointsParam = "false"
	}

	withPoints, err := strconv.ParseBool(withPointsParam)
	if err != nil {
		return terror.NewValidationError("bad validate", "withPoints")
	}

	userId := c.Get(middlewares.UserIdCtxClaim)
	ctx = context.WithValue(ctx, middlewares.UserIdCtxClaim, userId)

	response, err := h.service.GetRouteForTrip(ctx, tripId, withPoints)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
