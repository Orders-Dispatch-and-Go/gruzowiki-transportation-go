package carrierController

import (
	"github.com/gofiber/fiber/v2"
)

type CarrierController interface {
	GetCarrier(fiberCtx *fiber.Ctx) error
}
