package handlers

import (
	"net/http"
	"strconv"

	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"gruzowiki/services"

	"github.com/labstack/echo/v4"
)

type RecipientHandler struct {
	service *services.RecipientService
}

func NewRecipientHandler(service *services.RecipientService) *RecipientHandler {
	return &RecipientHandler{service: service}
}

func (h *RecipientHandler) CreateRecipient(c echo.Context) error {
	var req models.CreateRecipientRequest
	if err := c.Bind(&req); err != nil {
		return terror.NewValidationError(
			err.Error(),
			"binding CreateRecipientRequest",
		)
	}

	resp, err := h.service.CreateRecipient(
		c.Request().Context(),
		req.FirstName, req.SecondName, req.ThirdName,
		req.Phone, req.Email,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *RecipientHandler) GetRecipient(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return terror.NewValidationError(
			"invalid id",
			"parsing path parameter 'id'",
		)
	}

	resp, err := h.service.GetRecipient(c.Request().Context(), int32(id))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *RecipientHandler) ListRecipients(c echo.Context) error {
	resp, err := h.service.ListRecipients(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *RecipientHandler) UpdateRecipient(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return terror.NewValidationError(
			"invalid id",
			"parsing path parameter 'id'",
		)
	}

	var req models.UpdateRecipientRequest
	if err := c.Bind(&req); err != nil {
		return terror.NewValidationError(
			err.Error(),
			"binding UpdateRecipientRequest",
		)
	}

	resp, err := h.service.UpdateRecipient(c.Request().Context(), int32(id), req.FirstName, req.SecondName, req.ThirdName, req.Phone, req.Email)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *RecipientHandler) DeleteRecipient(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return terror.NewValidationError(
			"invalid id",
			"parsing path parameter 'id'",
		)
	}

	if err := h.service.DeleteRecipient(c.Request().Context(), int32(id)); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
