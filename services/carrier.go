package services

import (
	"context"
	"errors"
	"gruzowiki/db/pg"
	"gruzowiki/rest/exceptions"
	"gruzowiki/rest/models"
)

type Repo interface {
	GetCarrierById(context.Context, int32) (*pg.Carrier, error)
}

type CarrierService struct {
	repo Repo
}

func NewCarrierService(repo Repo) *CarrierService {
	return &CarrierService{
		repo: repo,
	}
}

func (c *CarrierService) GetCarrier(ctx context.Context, id int32) (*models.GetCarrierResponse, error) {
	carrier, err := c.repo.GetCarrierById(ctx, id)
	if err != nil {
		return nil, err
	}

	if carrier == nil {
		return nil, errors.New(exceptions.CarrierNotFound)
	}

	return &models.GetCarrierResponse{Id: carrier.ID, DriverCategory: carrier.DriverCategory.String}, err
}
