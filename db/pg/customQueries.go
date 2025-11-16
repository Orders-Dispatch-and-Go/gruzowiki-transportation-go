package pg

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func (q *Queries) RawQuery(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return q.db.Query(ctx, sql, args...)
}

func (q *Queries) RawQueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return q.db.QueryRow(ctx, sql, args...)
}

func (q *Queries) RawExec(ctx context.Context, sql string, args ...any) error {
	_, err := q.db.Exec(ctx, sql, args...)
	return err
}
