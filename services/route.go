package services

import (
	"context"
	"github.com/google/uuid"
	"gruzowiki/client"
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
		GetPointsForTrip(cargoRequestRouteID uuid.UUID) (*client.GetRoutePointsResponse, error)
		GetPointsForCargoRequest(cargoRequestRouteID uuid.UUID) (*client.GetRoutePointsResponse, error)
	}
)

func NewRouteService(client FeignClient, cargoRequestService *CargoRequestService, tripService *TripService) *RouteService {
	return &RouteService{
		client:              client,
		cargoRequestService: *cargoRequestService,
		tripService:         *tripService,
	}
}

func (s *RouteService) GetRouteForCargoRequest(ctx context.Context, cargoRequestId uuid.UUID, withPoints bool) (*models.GetTripRouteResponse, error) {
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

	if withPoints {
		points, err := s.client.GetPointsForCargoRequest(cargoRequestRouteId)
		if err != nil {
			return nil, err
		}

		route.Points = make([]float64, 0, len(points.Points)*2)
		for i := range points.Points {
			if len(points.Points[i]) != 2 {
				route.Points = nil
				break
			}
			route.Points = append(route.Points, points.Points[i][0], points.Points[i][1])
		}
	}

	route.ID = cargoRequestRouteId

	return route, nil
}

func (s *RouteService) GetRouteForTrip(ctx context.Context, tripId uuid.UUID, withPoints bool) (*models.GetTripRouteResponse, error) {
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

	if withPoints {
		points, err := s.client.GetPointsForTrip(tripRouteId)
		if err != nil {
			return nil, err
		}

		route.Points = make([]float64, 0, len(points.Points)*2)
		for i := range points.Points {
			if len(points.Points[i]) != 2 {
				route.Points = nil
				break
			}
			route.Points = append(route.Points, points.Points[i][0], points.Points[i][1])
		}
	}

	route.ID = tripRouteId

	return route, nil
}
