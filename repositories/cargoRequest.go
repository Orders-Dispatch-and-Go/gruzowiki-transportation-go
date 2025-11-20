package repositories

import (
	"context"
	"fmt"
	"gruzowiki/db/pg"
	"gruzowiki/rest/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"math/big"
)

type CargoRequestRepo struct {
	conn pg.Conn
}

func NewCargoRequestRepo(conn pg.Conn) *CargoRequestRepo {
	return &CargoRequestRepo{
		conn: conn,
	}
}

func (c *CargoRequestRepo) GetCargoRequestById(ctx context.Context, id uuid.UUID) (*pg.CargoRequest, error) {
	//cargoRequests, err := c.conn.Queries(ctx).RawQuery()
	return nil, nil
}

func (c *CargoRequestRepo) CreateCargoRequest(ctx context.Context, params pg.InsertCargoRequestParams) (*pgtype.UUID, error) {
	pgid, err := c.conn.Queries(ctx).InsertCargoRequest(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return &pgid, nil
}

func (c *CargoRequestRepo) GetCargoTypes(ctx context.Context) ([]pg.CargoType, error) {
	return c.conn.Queries(ctx).GetCargoTypes(ctx)
}

func (c *CargoRequestRepo) CreateCargo(ctx context.Context, cargos []models.Cargo) ([]pgtype.UUID, error) {
	ids := make([]pgtype.UUID, len(cargos))
	for i, cargo := range cargos {
		uuid, err := uuid.Parse(cargo.CargoRequestId)
		if err != nil {
			return nil, err
		}
		id, err := c.conn.Queries(ctx).CreateCargo(ctx, pg.CreateCargoParams{
			Length:    pgtype.Int4{Int32: int32(cargo.Length), Valid: true},
			Width:     pgtype.Int4{Int32: int32(cargo.Width), Valid: true},
			Height:    pgtype.Int4{Int32: int32(cargo.Height), Valid: true},
			Weight:    pgtype.Int4{Int32: int32(cargo.Weight), Valid: true},
			CargoType: pgtype.Int4{Int32: int32(cargo.CargoType), Valid: true},
			Worth:     pgtype.Numeric{Int: big.NewInt(int64(cargo.Worth)), Valid: true},
			RequestID: pgtype.UUID{Bytes: uuid, Valid: true},
		})
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}

	return ids, nil
}