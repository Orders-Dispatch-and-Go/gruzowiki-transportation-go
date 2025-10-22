package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gruzowiki/db/pg"
)

func NewConnect(ctx context.Context, dsn string) (pg.Conn, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("New(%q): %w", dsn, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Ping(): %w", err)
	}

	return pg.NewConn(pool), nil
}
