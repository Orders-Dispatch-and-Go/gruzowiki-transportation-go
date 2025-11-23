package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"gruzowiki/db/pg"
)

type StationRepo struct {
	conn pg.Conn
}

func NewStationRepo(conn pg.Conn) *StationRepo {
	return &StationRepo{
		conn: conn,
	}
}

func (r *StationRepo) InsertStation(ctx context.Context, params pg.InsertStationParams) (*pgtype.UUID, error) {
	pgid, err := r.conn.Queries(ctx).InsertStation(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return &pgid, nil
}

func (r *StationRepo) GetStations(ctx context.Context, ids []pgtype.UUID) ([]pg.Station, error) {
	stations, err := r.conn.Queries(ctx).SelectStations(ctx, pg.SelectStationsParams{
		ID:   ids[0],
		ID_2: ids[1],
	})

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return stations, nil
}
