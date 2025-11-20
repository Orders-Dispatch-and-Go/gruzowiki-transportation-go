package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"net/http"
	"strconv"
)

type CargoRequestHandler struct {
	service CargoRequestService
}

type CargoRequestService interface {
	SearchCargoRequests(
		ctx context.Context,
		filter models.GetCargoRequest,
		pageNumber int,
		pageSize int,
	) (*models.SearchCargoRequestsResponse, error)
	CreateCargoRequest(ctx context.Context, postCargoRequestRequest models.PostCargoRequestRequest) (*models.PostCargoRequestResponse, error)
}

func NewCargoRequestController(service CargoRequestService) *CargoRequestHandler {
	return &CargoRequestHandler{
		service: service,
	}
}

func (handler *CargoRequestHandler) GetCargoRequest(c echo.Context) error {
	ctx := c.Request().Context()

	pageNumber, err := strconv.Atoi(c.QueryParam("page_number"))
	if err != nil || pageNumber < 0 {
		return terror.NewValidationError("must be positive int and not null", "page_number")
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize < 1 {
		return terror.NewValidationError("must be positive int and not null", "page_size")
	}

	var request models.GetCargoRequest
	if err := c.Bind(&request); err != nil {
		return terror.NewValidationError("invalid request body", err.Error())
	}

	response, err := handler.service.SearchCargoRequests(ctx, request, pageNumber, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
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
