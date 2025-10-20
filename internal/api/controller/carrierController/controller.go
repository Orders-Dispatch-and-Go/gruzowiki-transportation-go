package carrierController

import (
	"auth-service/internal/api/controller"
	"auth-service/internal/service/carrierService"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Controller struct {
	requestReader controller.RequestReader
	service       carrierService.CarrierService
}

func New(requestReader controller.RequestReader, service carrierService.CarrierService) *Controller {
	return &Controller{requestReader, service}
}

func (c Controller) GetCarrier(fiberCtx *fiber.Ctx) error {
	id := fiberCtx.Params("id")

	carrierId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(err)
	}

	response, err := c.service.GetCarrier(int32(carrierId), fiberCtx.Context())

	if err != nil {
		// todo сделать глобальный общий обработчик ошибок
		return fiberCtx.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return fiberCtx.Status(fiber.StatusOK).JSON(response)
}
