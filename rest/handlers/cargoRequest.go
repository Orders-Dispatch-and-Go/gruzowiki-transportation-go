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
	GetCargoTypes(ctx context.Context) ([]models.CargoType, error)
	CreateCargo(ctx context.Context, cargo []models.Cargo) ([]string, error)
	MarkTrip(ctx context.Context, cargoReqId string, tripId string) error
}

func NewCargoRequestController(service CargoRequestService) *CargoRequestHandler {
	return &CargoRequestHandler{
		service: service,
	}
}

func (h *CargoRequestHandler) GetCargoRequest(c echo.Context) error {
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

	response, err := h.service.SearchCargoRequests(ctx, request, pageNumber, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *CargoRequestHandler) CreateCargoCargoRequest(c echo.Context) error {
	ctx := c.Request().Context()

	var request models.PostCargoRequestRequest
	if err := c.Bind(&request); err != nil {
		return terror.NewValidationError("invalid request body", err.Error())
	}

	response, err := h.service.CreateCargoRequest(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *CargoRequestHandler) GetCargoTypes(c echo.Context) error {
	ctx := c.Request().Context()
	cargoTypes, err := h.service.GetCargoTypes(ctx)
	if err != nil {
		return err
	}

	resp := models.CargoTypesResponse{CargoTypes: cargoTypes}

	return c.JSON(http.StatusOK, resp)
}

func (h *CargoRequestHandler) CreateCargo(c echo.Context) error {
	ctx := c.Request().Context()
	var req models.CargoRequest
	if err := c.Bind(&req); err != nil {
		return terror.NewValidationError("invalid request body", err.Error())
	}

	ids, err := h.service.CreateCargo(ctx, req.Cargo)
	if err != nil {
		return err
	}

	resp := models.IdsResponse{Ids: ids}

	return c.JSON(http.StatusOK, resp)
}

func (h* CargoRequestHandler) MarkTrip(c echo.Context) error {
	ctx := c.Request().Context()
	cargoReqId := c.Param("cargoRequestId")
	tripId := c.Param("tripId")

	if err := h.service.MarkTrip(ctx, cargoReqId, tripId); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}