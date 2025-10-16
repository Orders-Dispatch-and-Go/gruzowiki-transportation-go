package dto

import (
	"auth-service/internal/domain"
	repodto "auth-service/internal/repo/dto"
)

type GetUserRequest struct {
	UserID int64
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string
}

type RegisterRequest struct {
	Email     string
	Name      string
	Password  string
	Role      domain.UserRole
	Bio       *string
	AvatarURL *string
}

func (r RegisterRequest) ConvertToRepo(passwordHash string) repodto.InsertUserRequest {
	return repodto.InsertUserRequest{
		Email:        r.Email,
		Name:         r.Name,
		PasswordHash: passwordHash,
		Role:         r.Role,
		Bio:          r.Bio,
		AvatarURL:    r.AvatarURL,
	}
}

type VerifyEmailRequest struct {
	Token string
}

type SendVerificationRequest struct {
	UserID int64
}

func NewSendVerificationRequest(userID int64) SendVerificationRequest {
	return SendVerificationRequest{
		UserID: userID,
	}
}

func ConvertToRepoUserVerifiedRequest(userID int64, verified bool) repodto.SetUserVerifiedRequest {
	return repodto.SetUserVerifiedRequest{
		ID:       userID,
		Verified: verified,
	}
}

type ForgotPasswordRequest struct {
	Email string
}

type ResetPasswordRequest struct {
	Email    string
	Token    string
	Password string
}

func ConvertToRepoInsertPasswordResetTokenRequest(userID int64, token string) repodto.InsertPasswordResetTokenRequest {
	return repodto.InsertPasswordResetTokenRequest{
		ID:    userID,
		Token: token,
	}
}

func (r ResetPasswordRequest) ConvertToRepo(userID int64) repodto.GetPasswordResetTokenRequest {
	return repodto.GetPasswordResetTokenRequest{
		ID:    userID,
		Token: r.Token,
	}
}

func ConvertToRepoSetPasswordHashRequest(userID int64, passwordHash string) repodto.SetPasswordHashRequest {
	return repodto.SetPasswordHashRequest{
		ID:           userID,
		PasswordHash: passwordHash,
	}
}
