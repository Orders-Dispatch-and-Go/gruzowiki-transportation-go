package handlers

import (
	"gruzowiki/rest/models"
	"gruzowiki/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ConsignerHandler struct {
	service *services.ConsignerService
}

func NewConsignerHandler(service *services.ConsignerService) *ConsignerHandler {
	return &ConsignerHandler{service: service}
}

func (h *ConsignerHandler) CreateConsigner(c echo.Context) error {
	var req models.CreateConsignerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	err := h.service.CreateConsigner(c.Request().Context(), req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	resp := models.CreateConsignerResponse{ID: req.ID}
	return c.JSON(http.StatusOK, resp)
}