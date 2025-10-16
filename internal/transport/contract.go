package transport

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"auth-service/internal/domain"
	servicedto "auth-service/internal/service/dto"
)

type requestReader interface {
	ReadAndValidateFiberBody(fiberCtx *fiber.Ctx, request any) error
}

type service interface {
	GetUsers(ctx context.Context) ([]domain.User, error)
	GetUser(ctx context.Context, req servicedto.GetUserRequest) (domain.User, error)
	Login(ctx context.Context, req servicedto.LoginRequest) (servicedto.LoginResponse, error)
	Register(ctx context.Context, req servicedto.RegisterRequest) error
	VerifyEmail(ctx context.Context, req servicedto.VerifyEmailRequest) error
	SendVerification(ctx context.Context, req servicedto.SendVerificationRequest) error
	ForgotPassword(ctx context.Context, req servicedto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req servicedto.ResetPasswordRequest) error
}
