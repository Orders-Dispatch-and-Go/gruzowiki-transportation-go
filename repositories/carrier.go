package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type CarrierRepo struct {
	conn *sqlx.DB
}

func NewCarrierRepo(conn *sqlx.DB) *CarrierRepo {
	return &CarrierRepo{
		conn: conn,
	}
}

func (c *CarrierRepo) GetCarrierById(ctx context.Context, id string) (Carrier, error) {
	var carrier Carrier
	query := "SELECT id, driver_category FROM carrier WHERE id = $1"
	err := c.conn.GetContext(ctx, &carrier, query, id)
	if err != nil {
		return carrier, err
	}
	return carrier, nil
}
