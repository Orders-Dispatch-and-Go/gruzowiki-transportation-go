package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"net/http"
)

type CargoRequestHandler struct {
	service CargoRequestService
}

type CargoRequestService interface {
	CreateCargoRequest(ctx context.Context, postCargoRequestRequest models.PostCargoRequestRequest) (*models.PostCargoRequestResponse, error)
}

func NewCargoRequestController(service CargoRequestService) *CargoRequestHandler {
	return &CargoRequestHandler{
		service: service,
	}
}

func (handler *CargoRequestHandler) CreateCargoCargoRequest(c echo.Context) error {
	ctx := c.Request().Context()

	var request models.PostCargoRequestRequest
	if err := c.Bind(&request); err != nil {
		return terror.NewValidationError("invalid request body", err.Error())
	}

	response, err := handler.service.CreateCargoRequest(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response)
}
