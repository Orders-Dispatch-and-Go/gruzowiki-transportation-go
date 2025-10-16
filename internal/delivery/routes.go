package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func registerRoutes(app *fiber.App, transport transport, middlewares middlewares) {
	app.Use(recover.New())

	app.Post("/login", transport.Login)
	app.Post("/register", transport.Register)
	app.Get("/verify-email", transport.VerifyEmail)
	app.Post("/forgot-password", transport.ForgotPassword)
	app.Post("/reset-password", transport.ResetPassword)

	app.Use(middlewares.Auth)

	app.Get("/users", transport.GetUsers)
	app.Get("/user/:id", transport.GetUser)
	app.Get("/current-user", transport.GetCurrentUser)
	app.Post("/resend-verification", transport.ResendVerification)
}
