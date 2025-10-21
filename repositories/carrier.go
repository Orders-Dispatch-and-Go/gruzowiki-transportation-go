package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type CarrierRepo struct {
	conn *sqlx.DB
}

func NewCarrier(conn *sqlx.DB) (*CarrierRepo){
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