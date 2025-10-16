package delivery

import "github.com/gofiber/fiber/v2"

type middlewares interface {
	Auth(fiberCtx *fiber.Ctx) error
}

type transport interface {
	GetUser(fiberCtx *fiber.Ctx) error
	GetUsers(fiberCtx *fiber.Ctx) error
	GetCurrentUser(fiberCtx *fiber.Ctx) error
	Login(fiberCtx *fiber.Ctx) error
	Register(fiberCtx *fiber.Ctx) error
	VerifyEmail(fiberCtx *fiber.Ctx) error
	ResendVerification(fiberCtx *fiber.Ctx) error
	ForgotPassword(fiberCtx *fiber.Ctx) error
	ResetPassword(fiberCtx *fiber.Ctx) error
}
