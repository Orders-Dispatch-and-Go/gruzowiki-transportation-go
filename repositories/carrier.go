package repositories

import (
	"context"
	"errors"
	"gruzowiki/db/pg"
	"gruzowiki/terror"

	"github.com/jackc/pgx/v5"
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
			return nil, terror.NewObjectNotFound(id, "carrier")
		}
		return nil, err
	}

	return &carriers, nil
}
