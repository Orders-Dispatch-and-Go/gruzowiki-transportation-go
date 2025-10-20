package delivery

import (
	"auth-service/internal/api/controller/carrierController"
	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, carrierController carrierController.CarrierController, middlewares middlewares) {
	//app.Use(recover.New())

	//app.Use(middlewares.Auth)

	app.Get("/carriers/:id", carrierController.GetCarrier)
}
