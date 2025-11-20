package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"gruzowiki/db/pg"
	"gruzowiki/rest/models"
	"gruzowiki/util"
	"time"
)

type CargoRequestService struct {
	repo CargoRequestRepo
}

type CargoRequestRepo interface {
	CreateCargoRequest(ctx context.Context, params pg.InsertCargoRequestParams) (*pgtype.UUID, error)
	GetCargoTypes(ctx context.Context) ([]pg.CargoType, error)
	CreateCargo(ctx context.Context, cargos []models.Cargo) ([]pgtype.UUID, error)
}

func NewCargoRequestService(repo CargoRequestRepo) *CargoRequestService {
	return &CargoRequestService{
		repo: repo,
	}
}

func (c *CargoRequestService) CreateCargoRequest(ctx context.Context, postCargoRequestRequest models.PostCargoRequestRequest) (*models.PostCargoRequestResponse, error) {
	id, err := c.repo.CreateCargoRequest(ctx, pg.InsertCargoRequestParams{
		ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
		ConsignerID: pgtype.Int4{Int32: postCargoRequestRequest.ConsignerID, Valid: true},
		RecipientID: pgtype.Int4{Int32: postCargoRequestRequest.RecipientID, Valid: true},
		CreatedAt:   pgtype.Int8{Int64: time.Now().Unix(), Valid: true},
		Deadline:    pgtype.Int8{Int64: util.ToTimestamp(postCargoRequestRequest.Deadline), Valid: true},
		Price:       util.ToNumeric(postCargoRequestRequest.MaxPrice),
		Status:      pgtype.Text{String: models.StatusWaitingTripChoice, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	return &models.PostCargoRequestResponse{ID: id.Bytes}, nil
}

func (c *CargoRequestService) GetCargoTypes(ctx context.Context) ([]models.CargoType, error) {
	types, err := c.repo.GetCargoTypes(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]models.CargoType, len(types))
	for i, t := range types {
		resp[i] = models.CargoType{
			Id:      int(t.ID),
			Type:    t.Type.String,
			Fragile: t.Fragile.Bool,
		}
	}
	return resp, nil
}

func (c *CargoRequestService) CreateCargo(ctx context.Context, cargo []models.Cargo) ([]string, error) {
	ids, err := c.repo.CreateCargo(ctx, cargo)
	if err != nil {
		return nil, err
	}

	resp := make([]string, 0, len(ids))
	for i, _ := range ids {
		resp[i] = uuid.UUID(ids[i].Bytes).String()
	}

	return resp, nil
}