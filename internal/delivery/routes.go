package delivery

import (
	"auth-service/internal/api/controller/carrierController"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func registerRoutes(app *fiber.App, carrierController carrierController.CarrierController, middlewares middlewares) {
	app.Use(recover.New())

	app.Use(middlewares.Auth)

	app.Get("/carriers/:id", carrierController.GetCarrier)
}
