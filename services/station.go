package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"gruzowiki/db/pg"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"gruzowiki/util"
)

type StationServiceImpl struct {
	repo StationRepo
}

type StationRepo interface {
	InsertStation(ctx context.Context, params pg.InsertStationParams) (*pgtype.UUID, error)
	/* только для первых двух id */ GetStations(ctx context.Context, id []pgtype.UUID) ([]pg.Station, error)
}

func NewStationService(repo StationRepo) *StationServiceImpl {
	return &StationServiceImpl{
		repo: repo,
	}
}

func (s *StationServiceImpl) GetStations(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]models.Station, error) {
	if len(ids) == 0 {
		return nil, terror.NewValidationError("is empty", "ids")
	}

	pgUuids := make([]pgtype.UUID, 0, len(ids))
	for _, id := range ids {
		pgUuids = append(pgUuids, util.UuidToPgUuid(id))
	}

	pgStations, err := s.repo.GetStations(ctx, pgUuids)

	if err != nil {
		return nil, err
	}

	stationsMap := make(map[uuid.UUID]models.Station, len(pgStations))
	for _, station := range pgStations {
		stationsMap[util.PgUuidToUuid(station.ID)] = models.Station{
			Address: station.Address.String,
			Coords: models.Coords{
				Lat: station.Lat.Float64,
				Lon: station.Lon.Float64,
			},
		}
	}

	return stationsMap, nil
}

func (s *StationServiceImpl) CreateStation(ctx context.Context, station models.Station) (*models.CreateStationResponse, error) {
	pgid, err := s.repo.InsertStation(ctx, pg.InsertStationParams{
		ID:      util.UuidToPgUuid(uuid.New()),
		Address: util.GoTextToPgText(station.Address),
		Lat:     util.Float64ToPgFloat8(station.Coords.Lat),
		Lon:     util.Float64ToPgFloat8(station.Coords.Lon),
	})

	if err != nil {
		return nil, err
	}

	return &models.CreateStationResponse{ID: pgid.Bytes}, nil
}
