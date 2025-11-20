package repositories

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"gruzowiki/db/pg"
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
