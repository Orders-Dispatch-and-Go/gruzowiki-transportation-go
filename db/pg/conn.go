package pg

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn interface {
	Close()
	Queries(ctx context.Context) *Queries
	WithTx(ctx context.Context, txFunc func(ctx context.Context) error) error
}

type defaultConn struct {
	pool    *pgxpool.Pool
	queries *Queries
}

func NewConn(pool *pgxpool.Pool) Conn {
	return &defaultConn{
		pool:    pool,
		queries: New(pool),
	}
}

func (c *defaultConn) Close() {
	c.pool.Close()
}

func (c *defaultConn) Queries(ctx context.Context) *Queries {
	if tx := extractTx(ctx); tx != nil {
		return c.queries.WithTx(tx)
	}

	return c.queries
}

func (c *defaultConn) WithTx(ctx context.Context, txFunc func(ctx context.Context) error) error {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Begin(): %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("Rollback(): %v", err)
		}
	}()

	if err := txFunc(injectTx(ctx, tx)); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("Commit(): %w", err)
	}

	return nil
}
