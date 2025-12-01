package services

import (
	"context"
	"gruzowiki/repositories"
)

type ConsignerService struct {
	repo *repositories.ConsignerRepo
}

func NewConsignerService(repo *repositories.ConsignerRepo) *ConsignerService {
	return &ConsignerService{repo: repo}
}

func (s *ConsignerService) CreateConsigner(ctx context.Context, id int32) error {
	return s.repo.InsertConsigner(ctx, id)
}