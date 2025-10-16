package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"auth-service/internal/db/pg"
	"auth-service/internal/domain"
	"auth-service/internal/repo/dto"
)

type Repo struct {
	conn pg.Conn
}

func New(conn pg.Conn) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) GetUsers(ctx context.Context) ([]domain.User, error) {
	pgUsers, err := r.conn.Queries(ctx).GetUsers(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	users, err := dto.ConvertFromPgUsers(pgUsers)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return users, nil
}

func (r *Repo) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	pgUser, err := r.conn.Queries(ctx).GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // must be checked on the top level
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	user, err := dto.ConvertFromPgUser(pgUser)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return &user, nil
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	pgUser, err := r.conn.Queries(ctx).GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil //nolint:nilnil // must be checked on the top level
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	user, err := dto.ConvertFromPgUser(pgUser)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return &user, nil
}

func (r *Repo) InsertUser(ctx context.Context, req dto.InsertUserRequest) error {
	params, err := req.ConvertToParams()
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	err = r.conn.Queries(ctx).InsertUser(ctx, params)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}

func (r *Repo) SetUserVerified(ctx context.Context, req dto.SetUserVerifiedRequest) error {
	params := req.ConvertToParams()
	err := r.conn.Queries(ctx).SetUserVerified(ctx, params)

	return err
}

func (r *Repo) GetUserVerified(ctx context.Context, id int64) (bool, error) {
	verified, err := r.conn.Queries(ctx).GetUserVerificationStatus(ctx, id)
	return verified, err
}

func (r *Repo) InsertPasswordResetToken(ctx context.Context, req dto.InsertPasswordResetTokenRequest) error {
	params := req.ConvertToParams()
	err := r.conn.Queries(ctx).CreatePasswordResetToken(ctx, params)

	return err
}

func (r *Repo) GetPasswordResetToken(ctx context.Context, req dto.GetPasswordResetTokenRequest) (dto.GetPasswordResetTokenResponse, error) {
	params := req.ConvertToParams()

	pgPasswordResetToken, err := r.conn.Queries(ctx).GetPasswordResetToken(ctx, params)
	if err != nil {
		return dto.GetPasswordResetTokenResponse{}, fmt.Errorf("GetPasswordResetToken(%d, %s): %w", req.ID, req.Token, err)
	}

	passwordResetToken := dto.ConvertFromPgPasswordResetToken(pgPasswordResetToken)

	response := dto.NewGetPasswordResetTokenResponse(passwordResetToken)

	return response, nil
}

func (r *Repo) SetUserPasswordHash(ctx context.Context, req dto.SetPasswordHashRequest) error {
	params := req.ConvertToParams()
	err := r.conn.Queries(ctx).SetPasswordHash(ctx, params)

	return err
}

func (r *Repo) MarkPasswordResetTokenAsUsed(ctx context.Context, id int64) error {
	err := r.conn.Queries(ctx).MarkPasswordResetTokenAsUsed(ctx, id)
	return err
}
