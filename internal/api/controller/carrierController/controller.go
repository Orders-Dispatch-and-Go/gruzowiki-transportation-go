package carrierController

import (
	"auth-service/internal/api/controller"
	"auth-service/internal/service/carrierService"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	requestReader controller.RequestReader
	service       carrierService.CarrierService
}

func New(requestReader controller.RequestReader, service carrierService.CarrierService) *Controller {
	return &Controller{requestReader, service}
}

func (c Controller) GetCarrier(fiberCtx *fiber.Ctx) error {
	return nil
}
