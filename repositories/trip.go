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

func (r *TripRepo) GetTripsByIDsWithPagination(ctx context.Context, ids []string, limit, offset int32) ([]pg.GetTripsByIDsWithPaginationRow, error) {
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

func (r *TripRepo) CreateTrip(ctx context.Context, fromStation uuid.UUID, toStation uuid.UUID, startedAt int64, carrier int32) (uuid.UUID, error) {
	id, err := r.conn.Queries(ctx).InsertTrip(ctx, pg.InsertTripParams{
		FromStation: pgtype.UUID{Bytes: fromStation, Valid: true},
		ToStation:   pgtype.UUID{Bytes: toStation, Valid: true},
		StartedAt:   pgtype.Int8{Int64: startedAt, Valid: true},
		Carrier:     pgtype.Int4{Int32: carrier, Valid: true},
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("insert trip: %w", err)
	}

	return id.Bytes, nil
}
