package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gruzowiki/db/pg"
)

type CarrierRepo struct {
	conn pg.Conn
}

func NewCarrierRepo(conn pg.Conn) *CarrierRepo {
	return &CarrierRepo{
		conn: conn,
	}
}

func (c *CarrierRepo) GetCarrierById(ctx context.Context, id int32) (*pg.Carrier, error) {
	carriers, err := c.conn.Queries(ctx).GetCarrier(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return &carriers, nil
}
