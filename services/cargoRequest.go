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
