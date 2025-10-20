package middlewares

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"auth-service/internal/api/consts"
)

func (m *Middlewares) Auth(fiberCtx *fiber.Ctx) error {
	auth := fiberCtx.Get("Authorization")
	if auth == "" {
		return fiberCtx.Status(fiber.StatusUnauthorized).JSON(errorResponse{
			Error: "missing authorization header",
		})
	}

	authParts := strings.Split(auth, " ")
	if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
		return fiberCtx.Status(fiber.StatusUnauthorized).JSON(errorResponse{
			Error: "invalid authorization header format",
		})
	}

	token := authParts[1]

	userID, err := m.auth.ParseToken(token)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	fiberCtx.Locals(consts.UserIDContextKey, userID)

	return fiberCtx.Next()
}
