package repositories

import (
	"context"
	"errors"
	"fmt"
	"gruzowiki/db/pg"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CarrierRepo struct {
	conn pg.Conn
}

func NewCarrierRepo(conn pg.Conn) *CarrierRepo {
	return &CarrierRepo{conn: conn}
}

func (r *CarrierRepo) GetCarrierById(ctx context.Context, id int32) (*pg.Carrier, error) {
	carrier, err := r.conn.Queries(ctx).GetCarrier(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return &carrier, nil
}

func (r *CarrierRepo) CreateCarrier(ctx context.Context, id int32, driverCategory string) (int32, error) {
	pgCategory := pgtype.Text{
		String: driverCategory,
		Valid:  true,
	}

	params := pg.CreateCarrierParams{
		ID:             id,
		DriverCategory: pgCategory,
	}

	return r.conn.Queries(ctx).CreateCarrier(ctx, params)
}

func (r *CarrierRepo) UpdateCarrier(ctx context.Context, id int32, driverCategory string) error {
	_, err := r.conn.Queries(ctx).UpdateCarrier(ctx, pg.UpdateCarrierParams{
		ID: id,
		DriverCategory: pgtype.Text{
			String: driverCategory,
			Valid:  true,
		},
	})
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *CarrierRepo) DeleteCarrier(ctx context.Context, id int32) error {
	err := r.conn.Queries(ctx).DeleteCarrier(ctx, id)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}
