package api

import "github.com/gofiber/fiber/v2"

type requestReader interface {
	ReadAndValidateFiberBody(fiberCtx *fiber.Ctx, request any) error
}

type service interface {
}
