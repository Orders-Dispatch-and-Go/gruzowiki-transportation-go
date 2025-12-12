package handlers

import (
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
	var req models.CreateTripRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	id, err := h.service.CreateTrip(c.Request().Context(), req)
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
	tripIDParam := c.QueryParam("tripId")
	carrierIDParam := c.QueryParam("carrierId")

	if tripIDParam == "" && carrierIDParam == "" {
		return terror.NewValidationError("missing query parameters", "one of tripID or carrierId are required")
	}

	var tripID uuid.UUID
	if tripIDParam != "" {
		tripID2, err := uuid.Parse(tripIDParam)
		tripID = tripID2
		if err != nil {
			return terror.NewValidationError("invalid UUID", "tripID")
		}
	}

	var carrierID int
	if carrierIDParam != "" {
		carrierID2, err := strconv.Atoi(carrierIDParam)
		if err != nil {
			return terror.NewValidationError("invalid carrierId", "carrierId must be integer")
		}
		carrierID = carrierID2
	}

	trip, err := h.service.GetTripByIdAndCarrier(c.Request().Context(), tripID, int32(carrierID))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, trip)
}
