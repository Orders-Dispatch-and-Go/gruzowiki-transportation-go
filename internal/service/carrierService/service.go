package carrierService

import (
	"auth-service/internal/domain/carrierDto"
	"auth-service/internal/repo/carrierRepo"
	"context"
)

type Service struct {
	repo carrierRepo.CarrierRepo
}

func New(repo carrierRepo.CarrierRepo) *Service {
	return &Service{repo}
}

func (s *Service) GetCarrier(id int32, ctx context.Context) (
	getCarrierResponse carrierDto.GetCarrierResponse,
	err error,
) {
	return carrierDto.GetCarrierResponse{}, nil
}
