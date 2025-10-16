package domain

import "time"

type User struct {
	ID           int64
	Email        string
	Name         string
	PasswordHash string
	Role         UserRole
	Bio          *string
	AvatarURL    *string
	CreatedAt    *time.Time
	Verified     bool
}

type UserRole string

const (
	UserRoleReader = "reader"
	UserRoleWriter = "writer"
)
