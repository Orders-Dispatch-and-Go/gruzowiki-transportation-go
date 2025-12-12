package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"gruzowiki/db/pg"
	"gruzowiki/rest/terror"
)

type TripRepo struct {
	conn pg.Conn
}

func NewTripRepo(conn pg.Conn) *TripRepo {
	return &TripRepo{conn: conn}
}

func (r *TripRepo) GetTripById(ctx context.Context, tripId uuid.UUID) (*pg.Trip, error) {
	pgTrip, err := r.conn.Queries(ctx).GetTripById(ctx, pgtype.UUID{Bytes: tripId, Valid: true})
	if err != nil {
		return nil, err
	}
	if len(pgTrip) == 0 {
		return nil, err
	}
	return &pgTrip[0], err
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

func (r *TripRepo) CreateTrip(
	ctx context.Context,
	fromStation uuid.UUID,
	toStation uuid.UUID,
	routeId uuid.UUID,
	startedAt int64,
	carrier int32,
	carID *int32,
) (uuid.UUID, error) {

	var car pgtype.Int4
	if carID != nil {
		car = pgtype.Int4{Int32: *carID, Valid: true}
	} else {
		car = pgtype.Int4{Valid: false}
	}

	params := pg.InsertTripParams{
		FromStation:    pgtype.UUID{Bytes: fromStation, Valid: true},
		ToStation:      pgtype.UUID{Bytes: toStation, Valid: true},
		RouteID:        pgtype.UUID{Bytes: routeId, Valid: true},
		StartedAt:      pgtype.Int8{Int64: startedAt, Valid: true},
		CalculateEndAt: pgtype.Int8{Valid: false},
		ActualEndAt:    pgtype.Int8{Valid: false},
		Price:          pgtype.Numeric{Valid: false},
		Status:         pgtype.Text{String: "PENDING", Valid: true},
		Carrier:        pgtype.Int4{Int32: carrier, Valid: true},
		Car:            car,
	}

	id, err := r.conn.Queries(ctx).InsertTrip(ctx, params)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert trip: %w", err)
	}

	return id.Bytes, nil
}

func (r *TripRepo) UpdateTripStatus(ctx context.Context, tripID string, status string) error {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return terror.NewValidationError(err.Error(), "parsing path parameter 'id'")
	}
	return r.conn.Queries(ctx).UpdateTripStatus(ctx, pg.UpdateTripStatusParams{
		ID:     pgtype.UUID{Bytes: id, Valid: true},
		Status: pgtype.Text{String: status, Valid: true},
	})
}

func (r *TripRepo) UpdateRout(ctx context.Context, tripId, routeId string) error {
	id, err := uuid.Parse(tripId)
	if err != nil {
		terror.NewValidationError(err.Error(), "parsing path parameter 'id'")
	}
	route, err := uuid.Parse(routeId)
	if err != nil {
		terror.NewValidationError(err.Error(), "parsing path parameter 'id'")
	}
	return r.conn.Queries(ctx).UpdateTripRoute(ctx, pg.UpdateTripRouteParams{ID: pgtype.UUID{Bytes: id, Valid: true}, RouteID: pgtype.UUID{Bytes: route, Valid: true}})
}

func (r *TripRepo) GetTripByIdAndCarrier(ctx context.Context, tripID uuid.UUID, carrierID int32) (*pg.Trip, error) {
	pgID := pgtype.UUID{
		Bytes: tripID,
		Valid: true,
	}
	pgCarrier := pgtype.Int4{
		Int32: carrierID,
		Valid: true,
	}

	trip, err := r.conn.Queries(ctx).GetTripByIdAndCarrier(ctx, pg.GetTripByIdAndCarrierParams{
		ID:      pgID,
		Carrier: pgCarrier,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query GetTripByIdAndCarrier: %w", err)
	}

	return &trip, nil
}
