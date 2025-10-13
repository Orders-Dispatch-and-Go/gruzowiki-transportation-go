package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func (m *Middlewares) Authenticate(fiberCtx *fiber.Ctx) error {
	m.auth.ParseToken("")
	return fiberCtx.Next()
}
