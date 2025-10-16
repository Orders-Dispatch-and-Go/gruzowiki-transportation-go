package dto

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"auth-service/internal/db/pg"
	"auth-service/internal/domain"
)

func ConvertFromPgUsers(pgUsers []pg.User) ([]domain.User, error) {
	users := make([]domain.User, 0, len(pgUsers))

	for _, pgUser := range pgUsers {
		user, err := ConvertFromPgUser(pgUser)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func ConvertFromPgUser(pgUser pg.User) (domain.User, error) {
	role, err := ConvertFromPgUserRole(pgUser.Role)
	if err != nil {
		return domain.User{}, fmt.Errorf("role: %w", err)
	}

	user := domain.User{
		ID:           pgUser.ID,
		Email:        pgUser.Email,
		Name:         pgUser.Name,
		PasswordHash: pgUser.PasswordHash,
		Role:         role,
		Bio:          ConvertFromPgText(pgUser.Bio),
		AvatarURL:    ConvertFromPgText(pgUser.AvatarUrl),
		CreatedAt:    ConvertFromPgTimestamptz(pgUser.CreatedAt),
	}

	return user, nil
}

func ConvertFromPgUserRole(pgRole pg.UserRole) (domain.UserRole, error) {
	switch pgRole {
	case pg.UserRoleWriter:
		return domain.UserRoleWriter, nil
	case pg.UserRoleReader:
		return domain.UserRoleReader, nil
	default:
		return "", fmt.Errorf("unknown user role: %q", pgRole)
	}
}

func ConvertToPgUserRole(role domain.UserRole) (pg.UserRole, error) {
	switch role {
	case domain.UserRoleWriter:
		return pg.UserRoleWriter, nil
	case domain.UserRoleReader:
		return pg.UserRoleReader, nil
	default:
		return "", fmt.Errorf("unknown user role: %q", role)
	}
}

func ConvertFromPgText(text pgtype.Text) *string {
	if !text.Valid {
		return nil
	}

	return &text.String
}

func ConvertFromPgTimestamptz(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}

	return &ts.Time
}

func ConvertToPgText(str *string) pgtype.Text {
	if str == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{
		Valid:  true,
		String: *str,
	}
}
