package carrierService

import (
	"auth-service/internal/domain/carrierDto"
	"auth-service/internal/repo/carrierRepo"
	"context"
	"fmt"
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
	carrier, err := s.repo.GetCarrierById(id, ctx)

	if err != nil {
		// todo и вот это придется постоянно прописывать?!?!?!?!
		return carrierDto.GetCarrierResponse{}, err
	}

	if carrier == nil {
		return carrierDto.GetCarrierResponse{}, fmt.Errorf("carrier not found")
	}

	return carrierDto.GetCarrierResponse{
		ID:             carrier.ID,
		DriverCategory: carrier.DriverCategory,
	}, nil
}
