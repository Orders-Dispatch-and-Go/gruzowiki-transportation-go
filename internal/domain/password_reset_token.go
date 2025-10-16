package domain

import "time"

type PasswordResetToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiredAt *time.Time
	Used      bool
	CreatedAt *time.Time
}
