package dto

import (
	"fmt"

	"auth-service/internal/db/pg"
	"auth-service/internal/domain"
)

type InsertUserRequest struct {
	Email        string
	Name         string
	PasswordHash string
	Role         domain.UserRole
	Bio          *string
	AvatarURL    *string
}

func (r InsertUserRequest) ConvertToParams() (pg.InsertUserParams, error) {
	params := pg.InsertUserParams{
		Email:        r.Email,
		Name:         r.Name,
		PasswordHash: r.PasswordHash,
		Bio:          ConvertToPgText(r.Bio),
		AvatarUrl:    ConvertToPgText(r.AvatarURL),
	}

	var err error

	params.Role, err = ConvertToPgUserRole(r.Role)
	if err != nil {
		return params, fmt.Errorf("convert role: %w", err)
	}

	return params, nil
}

type SetUserVerifiedRequest struct {
	ID       int64
	Verified bool
}

func (r SetUserVerifiedRequest) ConvertToParams() pg.SetUserVerifiedParams {
	return pg.SetUserVerifiedParams{
		ID:       r.ID,
		Verified: r.Verified,
	}
}

type InsertPasswordResetTokenRequest struct {
	ID    int64
	Token string
}

func (r InsertPasswordResetTokenRequest) ConvertToParams() pg.CreatePasswordResetTokenParams {
	return pg.CreatePasswordResetTokenParams{
		UserID: r.ID,
		Token:  r.Token,
	}
}

type GetPasswordResetTokenRequest struct {
	ID    int64
	Token string
}

func (r GetPasswordResetTokenRequest) ConvertToParams() pg.GetPasswordResetTokenParams {
	return pg.GetPasswordResetTokenParams{
		UserID: r.ID,
		Token:  r.Token,
	}
}

type GetPasswordResetTokenResponse struct {
	Token domain.PasswordResetToken
}

func NewGetPasswordResetTokenResponse(token domain.PasswordResetToken) GetPasswordResetTokenResponse {
	return GetPasswordResetTokenResponse{
		Token: token,
	}
}

type SetPasswordHashRequest struct {
	ID           int64
	PasswordHash string
}

func (r SetPasswordHashRequest) ConvertToParams() pg.SetPasswordHashParams {
	return pg.SetPasswordHashParams{
		ID:           r.ID,
		PasswordHash: r.PasswordHash,
	}
}
