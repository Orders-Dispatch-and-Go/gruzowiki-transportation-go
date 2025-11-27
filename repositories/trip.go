package repositories

import (
	"context"
	"errors"
	"fmt"
	"gruzowiki/db/pg"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type TripRepo struct {
	conn pg.Conn
}

func NewTripRepo(conn pg.Conn) *TripRepo {
	return &TripRepo{conn: conn}
}

func (r *TripRepo) GetByCargoRequest(ctx context.Context, cargoRequestID uuid.UUID) (*pg.Trip, error) {
	pgID := pgtype.UUID{
		Bytes: cargoRequestID,
		Valid: true,
	}

	trip, err := r.conn.Queries(ctx).GetTripByCargoRequest(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query GetTripByCargoRequest: %w", err)
	}

	return &trip, nil
}