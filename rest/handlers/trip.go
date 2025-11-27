package handlers

import (
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"gruzowiki/services"
	"net/http"

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

	tripResp, err := h.service.GetTripByCargoRequest(c.Request().Context(), cargoRequestID)
	if err != nil {
		return err
	}

	resp := models.TripResponse{
		ID:              tripResp.ID,
		FromStation:     tripResp.FromStation,
		ToStation:       tripResp.ToStation,
		StartedAt:       tripResp.StartedAt,
		CalculatedEndAt: tripResp.CalculatedEndAt,
		ActualEndAt:     tripResp.ActualEndAt,
		Price:           tripResp.Price, 
		Status:          tripResp.Status,
		CarrierID:       tripResp.CarrierID,
		CarID:           tripResp.CarID,
	}

	return c.JSON(http.StatusOK, resp)
}