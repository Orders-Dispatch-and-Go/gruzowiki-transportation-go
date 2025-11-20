package repositories

import (
	"context"
	"errors"
	"fmt"
	"gruzowiki/db/pg"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CarRepo struct {
	conn pg.Conn
}

func NewCarRepo(conn pg.Conn) *CarRepo {
	return &CarRepo{conn: conn}
}

func (r *CarRepo) GetCar(ctx context.Context, id int32) (*pg.Car, error) {
	car, err := r.conn.Queries(ctx).GetCar(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return &car, nil
}

func (r *CarRepo) ListByOwner(ctx context.Context, ownerID int32) ([]pg.Car, error) {
	return r.conn.Queries(ctx).ListCarsByOwner(ctx, pgtype.Int4{
		Int32: ownerID,
		Valid: true,
	})
}

func (r *CarRepo) CreateCar(ctx context.Context, car pg.Car) (int32, error) {
	params := pg.CreateCarParams{
		Type:      car.Type,
		Length:    car.Length,
		Width:     car.Width,
		Height:    car.Height,
		MaxWeight: car.MaxWeight,
		Number:    car.Number,
		Owner:     car.Owner,
	}

	id, err := r.conn.Queries(ctx).CreateCar(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}

	return id, nil
}

func (r *CarRepo) UpdateCar(ctx context.Context, id int32, car pg.Car) error {
	params := pg.UpdateCarParams{
		ID:        id,
		Type:      car.Type,
		Length:    car.Length,
		Width:     car.Width,
		Height:    car.Height,
		MaxWeight: car.MaxWeight,
		Number:    car.Number,
		Owner:     car.Owner,
	}

	_, err := r.conn.Queries(ctx).UpdateCar(ctx, params)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}

func (r *CarRepo) DeleteCar(ctx context.Context, id int32) error {
	err := r.conn.Queries(ctx).DeleteCar(ctx, id)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}
