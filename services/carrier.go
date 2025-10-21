package services

import (
	"context"
	"gruzowiki/repositories"
)

type Repo interface {
	GetCarrierById(context.Context, string) (repositories.Carrier, error)
}

type CarrierService struct {
	repo Repo
}

func NewCarrierService(repo Repo) *CarrierService {
	return &CarrierService{
		repo: repo,
	}
}

func (c *CarrierService) GetCarrier(ctx context.Context, id string) (repositories.Carrier, error) {
	carrier, err := c.repo.GetCarrierById(ctx, id)
	if err != nil {
		return repositories.Carrier{}, nil
	}
	return carrier, err
}