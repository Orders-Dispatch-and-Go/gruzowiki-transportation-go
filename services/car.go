package services

import (
	"context"
	"gruzowiki/db/pg"
	"gruzowiki/repositories"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type CarService struct {
	carRepo     *repositories.CarRepo
	carrierRepo *repositories.CarrierRepo
}

func NewCarService(carRepo *repositories.CarRepo, carrierRepo *repositories.CarrierRepo) *CarService {
	return &CarService{
		carRepo:     carRepo,
		carrierRepo: carrierRepo,
	}
}

func (s *CarService) CreateCar(ctx context.Context, req models.CreateCarRequest) (*models.CreateCarResponse, error) {

	carrier, err := s.carrierRepo.GetCarrierById(ctx, req.OwnerID)
	if err != nil {
		return nil, err
	}
	if carrier == nil {
		return nil, terror.NewNotFoundError("Carrier", strconv.Itoa(int(req.OwnerID)))
	}

	car := pg.Car{
		Type:      pgtype.Text{String: req.Type, Valid: true},
		Length:    pgtype.Int4{Int32: req.Length, Valid: true},
		Width:     pgtype.Int4{Int32: req.Width, Valid: true},
		Height:    pgtype.Int4{Int32: req.Height, Valid: true},
		MaxWeight: pgtype.Int4{Int32: req.MaxWeight, Valid: true},
		Number:    pgtype.Text{String: req.Number, Valid: true},
		Owner:     pgtype.Int4{Int32: req.OwnerID, Valid: true},
	}

	id, err := s.carRepo.CreateCar(ctx, car)
	if err != nil {
		return nil, err
	}

	return &models.CreateCarResponse{ID: id}, nil
}

func (s *CarService) GetCar(ctx context.Context, id int32) (*models.GetCarResponse, error) {
	car, err := s.carRepo.GetCar(ctx, id)
	if err != nil {
		return nil, err
	}
	if car == nil {
		return nil, terror.NewNotFoundError("Car", strconv.Itoa(int(id)))
	}

	return &models.GetCarResponse{
		ID:        car.ID,
		Type:      car.Type.String,
		Length:    car.Length.Int32,
		Width:     car.Width.Int32,
		Height:    car.Height.Int32,
		MaxWeight: car.MaxWeight.Int32,
		Number:    car.Number.String,
		OwnerID:   car.Owner.Int32,
	}, nil
}

func (s *CarService) UpdateCar(ctx context.Context, id int32, req models.UpdateCarRequest) (*models.UpdateCarResponse, error) {
	existing, err := s.carRepo.GetCar(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, terror.NewNotFoundError("Car", strconv.Itoa(int(id)))
	}

	if req.OwnerID != nil {
		carrier, err := s.carrierRepo.GetCarrierById(ctx, *req.OwnerID)
		if err != nil {
			return nil, err
		}
		if carrier == nil {
			return nil, terror.NewNotFoundError("Carrier", strconv.Itoa(int(*req.OwnerID)))
		}
	}

	car := *existing

	if req.Type != nil {
		car.Type = pgtype.Text{String: *req.Type, Valid: true}
	}
	if req.Length != nil {
		car.Length = pgtype.Int4{Int32: *req.Length, Valid: true}
	}
	if req.Width != nil {
		car.Width = pgtype.Int4{Int32: *req.Width, Valid: true}
	}
	if req.Height != nil {
		car.Height = pgtype.Int4{Int32: *req.Height, Valid: true}
	}
	if req.MaxWeight != nil {
		car.MaxWeight = pgtype.Int4{Int32: *req.MaxWeight, Valid: true}
	}
	if req.Number != nil {
		car.Number = pgtype.Text{String: *req.Number, Valid: true}
	}
	if req.OwnerID != nil {
		car.Owner = pgtype.Int4{Int32: *req.OwnerID, Valid: true}
	}

	if err := s.carRepo.UpdateCar(ctx, id, car); err != nil {
		return nil, err
	}

	return &models.UpdateCarResponse{ID: id}, nil
}

func (s *CarService) DeleteCar(ctx context.Context, id int32) error {
	car, err := s.carRepo.GetCar(ctx, id)
	if err != nil {
		return err
	}
	if car == nil {
		return terror.NewNotFoundError("Car", strconv.Itoa(int(id)))
	}

	return s.carRepo.DeleteCar(ctx, id)
}

func (s *CarService) ListCarsByOwner(ctx context.Context, ownerID int32) ([]*models.GetCarResponse, error) {
	carrier, err := s.carrierRepo.GetCarrierById(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	if carrier == nil {
		return nil, terror.NewNotFoundError("Carrier", strconv.Itoa(int(ownerID)))
	}

	cars, err := s.carRepo.ListByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	res := make([]*models.GetCarResponse, len(cars))
	for i, car := range cars {
		res[i] = &models.GetCarResponse{
			ID:        car.ID,
			Type:      car.Type.String,
			Length:    car.Length.Int32,
			Width:     car.Width.Int32,
			Height:    car.Height.Int32,
			MaxWeight: car.MaxWeight.Int32,
			Number:    car.Number.String,
			OwnerID:   car.Owner.Int32,
		}
	}

	return res, nil
}
