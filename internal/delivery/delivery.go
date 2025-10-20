package delivery

import (
	"auth-service/internal/api/controller/carrierController"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Delivery struct {
	cfg Config
	app *fiber.App
}

func New(cfg Config, carrierController carrierController.CarrierController, middlewares middlewares) *Delivery {
	app := fiber.New(cfg.Server.convertToForeign())
	app.Use(cors.New(cfg.Cors.convertToForeign()))

	registerRoutes(app, carrierController, middlewares)

	return &Delivery{
		cfg: cfg,
		app: app,
	}
}

func (d *Delivery) Listen() error {
	return d.app.Listen(d.cfg.Serve.Address)
}

func (d *Delivery) Shutdown(ctx context.Context) error {
	return d.app.ShutdownWithContext(ctx)
}
