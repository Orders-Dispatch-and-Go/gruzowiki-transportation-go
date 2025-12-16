package handlers

import (
	"context"
	"gruzowiki/rest/middlewares"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"gruzowiki/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TripHandler struct {
	service *services.TripService
}

func NewTripHandler(service *services.TripService) *TripHandler {
	return &TripHandler{service: service}
}

func (h *TripHandler) GetTripByCargoRequest(c echo.Context) error {
	idParam := c.Param("id")
	cargoRequestID, err := uuid.Parse(idParam)
	if err != nil {
		return terror.NewValidationError("invalid UUID", "parsing path parameter 'id'")
	}

	pageNumber := c.QueryParam("page_number")
	pageSize := c.QueryParam("page_size")

	if pageNumber != "" && pageSize == "" {
		pageNum, err := strconv.Atoi(pageNumber)
		if err != nil || pageNum < 1 {
			pageNum = 1
		}

		pageSz, err := strconv.Atoi(pageSize)
		if err != nil || pageSz < 1 {
			pageSz = 10
		}

		tripResp, err := h.service.GetTripsByCargoRequest(c.Request().Context(), cargoRequestID, pageNum, pageSz)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, tripResp)
	}

	tripResp, err := h.service.GetTripByCargoRequest(c.Request().Context(), cargoRequestID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, *tripResp)
}

func (h *TripHandler) CreateTrip(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateTripRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	userId := c.Get(middlewares.UserIdCtxClaim)
	ctx = context.WithValue(ctx, middlewares.UserIdCtxClaim, userId)

	id, err := h.service.CreateTrip(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": id.String()})
}

func (h *TripHandler) Finish(c echo.Context) error {
	id := c.Param("id")
	status := c.Param("status")

	if err := h.service.FinishTrip(c.Request().Context(), id, status); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *TripHandler) Start(c echo.Context) error {
	id := c.Param("id")
	var req models.CargoRequestIdsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	err := h.service.StartTrip(c.Request().Context(), id, req.Ids)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *TripHandler) GetTripByIdAndCarrier(c echo.Context) error {
	ctx := c.Request().Context()

	tripIDParam := c.QueryParam("tripId")
	carrierIDParam := c.QueryParam("carrierId")

	var tripID *uuid.UUID
	if tripIDParam != "" {
		id, err := uuid.Parse(tripIDParam)
		if err != nil {
			return terror.NewValidationError("invalid tripId", "must be uuid")
		}
		tripID = &id
	}

	var carrierID *int32
	if carrierIDParam != "" {
		id, err := strconv.Atoi(carrierIDParam)
		if err != nil {
			return terror.NewValidationError("invalid carrierId", "must be integer")
		}
		cid := int32(id)
		carrierID = &cid
	}

	trip, err := h.service.GetTripByIdAndCarrier(ctx, tripID, carrierID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, trip)
}
