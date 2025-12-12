package services

import (
	"context"
	"github.com/google/uuid"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
)

type (
	RouteService struct {
		client              FeignClient
		cargoRequestService CargoRequestService
		tripService         TripService
	}

	FeignClient interface {
		GetRouteForCargoRequest(cargoRequestRouteID uuid.UUID) (*models.GetTripRouteResponse, error)
		GetRouteForTrip(tripRouteId uuid.UUID) (*models.GetTripRouteResponse, error)
	}
)

func NewRouteService(client FeignClient, cargoRequestService *CargoRequestService, tripService *TripService) *RouteService {
	return &RouteService{
		client:              client,
		cargoRequestService: *cargoRequestService,
		tripService:         *tripService,
	}
}

func (s *RouteService) GetRouteForCargoRequest(ctx context.Context, cargoRequestId uuid.UUID) (*models.GetTripRouteResponse, error) {
	stringId := cargoRequestId.String()
	response, err := s.cargoRequestService.SearchCargoRequests(
		ctx,
		models.GetCargoRequest{
			ID:          &stringId,
			ConsignerID: nil,
			RecipientID: nil,
			Status:      nil,
			CreatedFrom: nil,
			CreatedTo:   nil,
		},
		1,
		10,
	)
	if err != nil {
		return nil, err
	}

	if len(response.CargoRequests) == 0 {
		return nil, terror.NewNotFoundError("cargoRequest", "uuid")
	}

	if response.CargoRequests[0].RouteID == nil {
		return nil, terror.NewNotFoundError("route", "cargoRequest")
	}

	cargoRequestRouteId, err := uuid.Parse(*response.CargoRequests[0].RouteID)
	if err != nil {
		return nil, terror.NewValidationError("invalid UUID", "parsing path parameter 'uuid'")
	}

	route, err := s.client.GetRouteForCargoRequest(cargoRequestRouteId)
	if err != nil {
		return nil, err
	}

	route.ID = cargoRequestRouteId

	return route, nil
}

func (s *RouteService) GetRouteForTrip(ctx context.Context, tripId uuid.UUID) (*models.GetTripRouteResponse, error) {
	response, err := s.tripService.GetTripById(ctx, tripId)
	if err != nil {
		return nil, err
	}

	if response.RouteID == nil {
		return nil, terror.NewNotFoundError("route", "trip")
	}

	uuidRouteId := *response.RouteID
	tripRouteId, err := uuid.Parse(uuidRouteId.String())
	if err != nil {
		return nil, terror.NewValidationError("invalid UUID", "parsing path parameter 'uuid'")
	}

	route, err := s.client.GetRouteForTrip(tripRouteId)
	if err != nil {
		return nil, err
	}

	route.ID = tripRouteId

	return route, nil
}
