package carrierRepo

import (
	"auth-service/internal/db/pg"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	conn pg.Conn
}

func New(conn pg.Conn) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) GetCarrierById(id int32, ctx context.Context) (*pg.Carrier, error) {
	carriers, err := r.conn.Queries(ctx).GetCarrier(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return &carriers, nil
}
