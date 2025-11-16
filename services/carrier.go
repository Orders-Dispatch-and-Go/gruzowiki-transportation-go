package services

import (
	"context"
	"gruzowiki/repositories"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"strconv"
)

type CarrierService struct {
	repo *repositories.CarrierRepo
}

func NewCarrierService(repo *repositories.CarrierRepo) *CarrierService {
	return &CarrierService{repo: repo}
}

func (s *CarrierService) CreateCarrier(ctx context.Context, id int32, driverCategory string) (*models.CreateCarrierResponse, error) {
	newID, err := s.repo.CreateCarrier(ctx, id, driverCategory)
	if err != nil {
		return nil, err
	}
	return &models.CreateCarrierResponse{ID: newID}, nil
}

func (s *CarrierService) GetCarrier(ctx context.Context, id int32) (*models.GetCarrierResponse, error) {
	carrier, err := s.repo.GetCarrierById(ctx, id)
	if err != nil {
		return nil, err
	}
	if carrier == nil {
		return nil, terror.NewNotFoundError("Carrier", strconv.Itoa(int(id)))
	}

	driverCategory := ""
	if carrier.DriverCategory.Valid {
		driverCategory = carrier.DriverCategory.String
	}

	return &models.GetCarrierResponse{
		ID:             carrier.ID,
		DriverCategory: driverCategory,
	}, nil
}

func (s *CarrierService) UpdateCarrier(ctx context.Context, id int32, driverCategory string) (*models.UpdateCarrierResponse, error) {
	c, err := s.repo.GetCarrierById(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, terror.NewNotFoundError("Carrier", strconv.Itoa(int(id)))
	}

	err = s.repo.UpdateCarrier(ctx, id, driverCategory)
	if err != nil {
		return nil, err
	}
	return &models.UpdateCarrierResponse{ID: id}, nil
}

func (s *CarrierService) DeleteCarrier(ctx context.Context, id int32) error {
	c, err := s.repo.GetCarrierById(ctx, id)
	if err != nil {
		return err
	}
	if c == nil {
		return terror.NewNotFoundError("Carrier", strconv.Itoa(int(id)))
	}

	return s.repo.DeleteCarrier(ctx, id)
}
