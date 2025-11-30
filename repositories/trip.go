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

func (r *TripRepo) GetTripsIdByCargoRequest(ctx context.Context, cargoRequestID uuid.UUID) ([]pgtype.UUID, error) {
	pgID := pgtype.UUID{
		Bytes: cargoRequestID,
		Valid: true,
	}
	ids, err := r.conn.Queries(ctx).GetSuitableTripsForCargoRequest(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query GetTripByCargoRequest: %w", err)
	}
	return ids, nil
}

func (r *TripRepo) GetTripsByIDsWithPagination(ctx context.Context, ids []string,  limit, offset int32) ([]pg.GetTripsByIDsWithPaginationRow, error) {
	uuids := make([]pgtype.UUID, len(ids))
	for i, id := range ids {
		uuid, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		uuids[i] = pgtype.UUID{
			Bytes: uuid,
			Valid: true,
		}
	}
	return r.conn.Queries(ctx).GetTripsByIDsWithPagination(ctx, pg.GetTripsByIDsWithPaginationParams{Column1: uuids, Limit: limit, Offset: offset})
}