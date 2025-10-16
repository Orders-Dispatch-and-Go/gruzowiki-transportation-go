package transport

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"auth-service/internal/transport/consts"
)

func extractUserID(fiberCtx *fiber.Ctx) (int64, error) {
	userID, ok := fiberCtx.Locals(consts.UserIDContextKey).(int64)
	if !ok {
		return 0, errors.New("no user ID in context")
	}

	return userID, nil
}
