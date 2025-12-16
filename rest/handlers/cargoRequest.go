package handlers

import (
	"context"
	"gruzowiki/rest/middlewares"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
	Delivered(ctx context.Context, cargoReqId string, code string) error
	GetRequestsForTripWithRoutes(
		ctx context.Context,
		tripID uuid.UUID,
		filter models.GetCargoRequestsForTripFilter,
	) (*models.GetCargoRequestsForTripResponse, error)
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

	userId := c.Get(middlewares.UserIdCtxClaim)
	ctx = context.WithValue(ctx, middlewares.UserIdCtxClaim, userId)

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

	userId := c.Get(middlewares.UserIdCtxClaim)
	ctx = context.WithValue(ctx, middlewares.UserIdCtxClaim, userId)

	if request.ConsignerID != userId {
		return terror.NewValidationError("must be equal to userId from jwt", "consignerId")
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

func (h *CargoRequestHandler) GetRequestsForTrip(c echo.Context) error {
	ctx := c.Request().Context()

	tripIDStr := c.Param("tripID")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		return terror.NewValidationError("invalid tripID", tripIDStr)
	}

	//userId := c.Get(middlewares.UserIdCtxClaim)
	//ctx = context.WithValue(ctx, middlewares.UserIdCtxClaim, userId)

	var filter models.GetCargoRequestsForTripFilter

	if val := c.QueryParam("cargoLengthMax"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			tmp := int32(i)
			filter.CargoLengthMax = &tmp
		} else {
			return terror.NewValidationError("invalid cargoLengthMax", val)
		}
	}

	if val := c.QueryParam("cargoWidthMax"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			tmp := int32(i)
			filter.CargoWidthMax = &tmp
		} else {
			return terror.NewValidationError("invalid cargoWidthMax", val)
		}
	}

	if val := c.QueryParam("cargoHeightMax"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			tmp := int32(i)
			filter.CargoHeightMax = &tmp
		} else {
			return terror.NewValidationError("invalid cargoHeightMax", val)
		}
	}

	if val := c.QueryParam("cargoType"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			tmp := int32(i)
			filter.CargoType = &tmp
		} else {
			return terror.NewValidationError("invalid cargoType", val)
		}
	}

	if val := c.QueryParam("deadline"); val != "" {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			ts := t.Unix()
			filter.Deadline = &ts
		} else {
			return terror.NewValidationError("invalid deadline", val)
		}
	}

	if val := c.QueryParam("minPrice"); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			tmp := i
			filter.MinPrice = &tmp
		} else {
			return terror.NewValidationError("invalid minPrice", val)
		}
	}

	resp, err := h.service.GetRequestsForTripWithRoutes(ctx, tripID, filter)
	if err != nil {
		return err
	}

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

func (h *CargoRequestHandler) MarkTrip(c echo.Context) error {
	ctx := c.Request().Context()
	cargoReqId := c.Param("cargoRequestId")
	tripId := c.Param("tripId")

	if err := h.service.MarkTrip(ctx, cargoReqId, tripId); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *CargoRequestHandler) Delivered(c echo.Context) error {
	ctx := c.Request().Context()
	cargoReqId := c.Param("uuid")
	code := c.Param("code")

	if err := h.service.Delivered(ctx, cargoReqId, code); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
