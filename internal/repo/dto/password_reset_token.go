package dto

import (
	"auth-service/internal/db/pg"
	"auth-service/internal/domain"
)

func ConvertFromPgPasswordResetToken(pgPasswordResetToken pg.PasswordResetToken) domain.PasswordResetToken {
	return domain.PasswordResetToken{
		ID:        pgPasswordResetToken.ID,
		UserID:    pgPasswordResetToken.UserID,
		Token:     pgPasswordResetToken.Token,
		ExpiredAt: ConvertFromPgTimestamptz(pgPasswordResetToken.ExpiredAt),
		Used:      pgPasswordResetToken.Used,
		CreatedAt: ConvertFromPgTimestamptz(pgPasswordResetToken.CreatedAt),
	}
}
