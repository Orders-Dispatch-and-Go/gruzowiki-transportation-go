package delivery

import "github.com/gofiber/fiber/v2"

type middlewares interface {
	Auth(fiberCtx *fiber.Ctx) error
}
