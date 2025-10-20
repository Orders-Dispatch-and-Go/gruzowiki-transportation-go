package controller

import "github.com/gofiber/fiber/v2"

type RequestReader interface {
	ReadAndValidateFiberBody(fiberCtx *fiber.Ctx, request any) error
}
