package repositories

import (
	"context"
	"fmt"
	"gruzowiki/db/pg"
)

type ConsignerRepo struct {
	conn pg.Conn
}

func NewConsignerRepo(conn pg.Conn) *ConsignerRepo {
	return &ConsignerRepo{conn: conn}
}

func (r *ConsignerRepo) InsertConsigner(ctx context.Context, id int32) error {
	err := r.conn.Queries(ctx).InsertConsigner(ctx, id)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	return nil
}
