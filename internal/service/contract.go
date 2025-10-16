package service

import (
	"context"

	"auth-service/internal/domain"
	"auth-service/internal/repo/dto"
	repodto "auth-service/internal/repo/dto"
)

type auth interface {
	ParseToken(tokenString string) (int64, error)
	GenerateToken(userID int64) (string, error)
	GenerateVerificationToken(userID int64) (string, error)
}

type passwordHasher interface {
	Hash(password string) (string, error)
	VerifyHash(password, passwordHash string) error
}

type repo interface {
	GetUsers(ctx context.Context) ([]domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	InsertUser(ctx context.Context, req repodto.InsertUserRequest) error
	SetUserVerified(ctx context.Context, req repodto.SetUserVerifiedRequest) error
	GetUserVerified(ctx context.Context, id int64) (bool, error)
	InsertPasswordResetToken(ctx context.Context, req repodto.InsertPasswordResetTokenRequest) error
	GetPasswordResetToken(ctx context.Context, req dto.GetPasswordResetTokenRequest) (dto.GetPasswordResetTokenResponse, error)
	SetUserPasswordHash(ctx context.Context, req dto.SetPasswordHashRequest) error
	MarkPasswordResetTokenAsUsed(ctx context.Context, id int64) error
}

type emailSender interface {
	SendVerificationEmail(email, token string) error
	SendPasswordResetEmail(email, token string) error
}
