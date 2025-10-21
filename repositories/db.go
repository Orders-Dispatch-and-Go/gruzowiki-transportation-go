package repositories

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewConnect(address string) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("pgx", address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}