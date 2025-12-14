package services

import (
	"context"
	"fmt"
	"gruzowiki/db/pg"
	"gruzowiki/repositories"
	"gruzowiki/rest/models"
	"gruzowiki/rest/terror"
	"gruzowiki/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type TripFeignClient interface {
	GetPotentialTrips(cargoRequestRouteID string, tripRouteIDs []string) ([]string, error)
	CreateRouteForTrip(
		request models.CreateTripRequest,
		fromStationId uuid.UUID,
		toStationId uuid.UUID,
	) (*uuid.UUID, error)
	MergeRoutes(tripRouteID string, cargoRequestRouteIDs []string) (uuid.UUID, error)
}

type TripService struct {
	tripRepo         *repositories.TripRepo
	stationRepo      *repositories.StationRepo
	cargoRequestRepo *repositories.CargoRequestRepo
	client           TripFeignClient
}

func NewTripService(
	tripRepo *repositories.TripRepo,
	stationRepo *repositories.StationRepo,
	cargoRequestRepo *repositories.CargoRequestRepo,
	client TripFeignClient,
) *TripService {
	return &TripService{
		tripRepo:         tripRepo,
		stationRepo:      stationRepo,
		client:           client,
		cargoRequestRepo: cargoRequestRepo,
	}
}

func (s *TripService) GetTripById(ctx context.Context, tripId uuid.UUID) (models.TripResponse, error) {
	pgTrip, err := s.tripRepo.GetTripById(ctx, tripId)
	if err != nil {
		return models.TripResponse{}, err
	}

	if pgTrip == nil {
		return models.TripResponse{}, terror.NewNotFoundError("trip", tripId.String())
	}

	trip := models.TripResponse{
		ID:          pgTrip.ID.Bytes,
		RouteID:     util.PgUuidToGoUuidPointer(pgTrip.RouteID),
		FromStation: models.Station{},
		ToStation:   models.Station{},
		//StartedAt:   pgTrip.StartedAt.Int64,
		ActualEndAt: pgTrip.ActualEndAt.Int64,
		Price:       util.NumericToString(pgTrip.Price),
		Status:      pgTrip.Status.String,
		CarrierID:   int(pgTrip.Carrier.Int32),
		CarID:       int(pgTrip.Car.Int32),
	}

	if util.PgUuidToGoUuidPointer(pgTrip.FromStation) != nil && util.PgUuidToGoUuidPointer(pgTrip.ToStation) != nil {
		stations, err := s.stationRepo.GetStations(ctx, []pgtype.UUID{
			pgTrip.FromStation,
			pgTrip.ToStation,
		})
		if err != nil {
			return models.TripResponse{}, err
		}

		var fromStation, toStation models.Station
		for _, st := range stations {
			id := uuid.UUID(st.ID.Bytes)
			if id == pgTrip.FromStation.Bytes {
				fromStation = models.Station{
					Address: st.Address.String,
					Coords:  models.Coords{Lat: st.Lat.Float64, Lon: st.Lon.Float64},
				}
			}
			if id == pgTrip.ToStation.Bytes {
				toStation = models.Station{
					Address: st.Address.String,
					Coords:  models.Coords{Lat: st.Lat.Float64, Lon: st.Lon.Float64},
				}
			}
		}

		trip.FromStation = fromStation
		trip.ToStation = toStation
	}

	return trip, nil
}

func (s *TripService) GetTripByCargoRequest(ctx context.Context, cargoRequestID uuid.UUID) (*models.TripResponse, error) {
	trip, err := s.tripRepo.GetByCargoRequest(ctx, cargoRequestID)
	if err != nil {
		return nil, err
	}
	if trip == nil {
		return nil, terror.NewNotFoundError("Trip", cargoRequestID.String())
	}

	fromID := pgtype.UUID{Bytes: trip.FromStation.Bytes, Valid: true}
	toID := pgtype.UUID{Bytes: trip.ToStation.Bytes, Valid: true}

	stations, err := s.stationRepo.GetStations(ctx, []pgtype.UUID{fromID, toID})
	if err != nil {
		return nil, err
	}

	var fromStation, toStation models.Station
	for _, st := range stations {
		id := uuid.UUID(st.ID.Bytes)
		if id == trip.FromStation.Bytes {
			fromStation = models.Station{
				Address: st.Address.String,
				Coords:  models.Coords{Lat: st.Lat.Float64, Lon: st.Lon.Float64},
			}
		}
		if id == trip.ToStation.Bytes {
			toStation = models.Station{
				Address: st.Address.String,
				Coords:  models.Coords{Lat: st.Lat.Float64, Lon: st.Lon.Float64},
			}
		}
	}

	return &models.TripResponse{
		ID:              trip.ID.Bytes,
		FromStation:     fromStation,
		ToStation:       toStation,
		StartedAt:       trip.StartedAt.Int64,
		CalculatedEndAt: trip.CalculateEndAt.Int64,
		ActualEndAt:     trip.ActualEndAt.Int64,
		Price:           util.NumericToString(trip.Price),
		Status:          trip.Status.String,
		CarrierID:       int(trip.Carrier.Int32),
		CarID:           int(trip.Car.Int32),
	}, nil
}

func (s *TripService) GetTripsByCargoRequest(ctx context.Context, cargoRequestID uuid.UUID, pageNumber, pageSize int) (models.Trips, error) {
	ids, err := s.tripRepo.GetTripsIdByCargoRequest(ctx, cargoRequestID)
	if err != nil {
		return models.Trips{}, err
	}

	suitsTripsIds, err := s.client.GetPotentialTrips(cargoRequestID.String(), uuidsToStrigs(ids))
	if err != nil {
		return models.Trips{}, err
	}

	trips, err := s.tripRepo.GetTripsByIDsWithPagination(ctx, suitsTripsIds, int32(pageSize), int32(pageNumber)*int32(pageSize))
	if err != nil {
		return models.Trips{}, err
	}

	resp := models.Trips{Trips: make([]models.TripResponse, 0, 0)}
	for _, trip := range trips {
		resp.Trips = append(resp.Trips, models.TripResponse{
			ID: trip.ID.Bytes,
			FromStation: models.Station{
				Address: trip.FromAddress.String,
				Coords:  models.Coords{Lat: trip.FromLat.Float64, Lon: trip.FromLon.Float64},
			},
			ToStation: models.Station{
				Address: trip.ToAddress.String,
				Coords:  models.Coords{Lat: trip.ToLat.Float64, Lon: trip.ToLon.Float64},
			},
			StartedAt:       trip.StartedAt.Int64,
			CalculatedEndAt: trip.CalculateEndAt.Int64,
			ActualEndAt:     trip.ActualEndAt.Int64,
			Price:           util.NumericToString(trip.Price),
			Status:          trip.Status.String,
			CarrierID:       int(trip.CarID.Int32),
			CarID:           int(trip.CarID.Int32),
		})
	}

	return resp, nil
}

func uuidsToStrigs(ids []pgtype.UUID) []string {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = uuid.UUID(id.Bytes).String()
	}
	return strs
}

func (s *TripService) CreateTrip(ctx context.Context, req models.CreateTripRequest) (*uuid.UUID, error) {

	fromID, err := s.stationRepo.InsertStation(ctx, pg.InsertStationParams{
		ID:      pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Address: pgtype.Text{String: req.FromStation.Address, Valid: true},
		Lat:     pgtype.Float8{Float64: req.FromStation.Coords.Lat, Valid: true},
		Lon:     pgtype.Float8{Float64: req.FromStation.Coords.Lon, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("insert from station: %w", err)
	}

	toID, err := s.stationRepo.InsertStation(ctx, pg.InsertStationParams{
		ID:      pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Address: pgtype.Text{String: req.ToStation.Address, Valid: true},
		Lat:     pgtype.Float8{Float64: req.ToStation.Coords.Lat, Valid: true},
		Lon:     pgtype.Float8{Float64: req.ToStation.Coords.Lon, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("insert to station: %w", err)
	}

	routeID, err := s.client.CreateRouteForTrip(req, fromID.Bytes, toID.Bytes)
	if err != nil {
		return nil, fmt.Errorf("create route for trip: %w", err)
	}

	startedAtInt := util.ToTimestamp(req.StartedAt)

	tripID, err := s.tripRepo.CreateTrip(
		ctx,
		fromID.Bytes,
		toID.Bytes,
		*routeID,
		startedAtInt,
		req.Carrier,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("insert trip: %w", err)
	}

	return &tripID, nil
}

func (s *TripService) FinishTrip(ctx context.Context, tripID string, status string) error {
	if status != models.TripStatusCompleted && status != models.TripStatusCanceled {
		return terror.NewValidationError("you cannot complete the trip with the transferred status", "status")
	}
	return s.tripRepo.UpdateTripStatus(ctx, tripID, status)
}

func (s *TripService) StartTrip(ctx context.Context, tripId string, cargoRequestIds []string) error {
	cargoRequestRoutes := make([]pg.CargoRequest, 0, len(cargoRequestIds))
	for _, cargoRequestId := range cargoRequestIds {
		cargoRequestUUID, err := uuid.Parse(cargoRequestId)
		if err != nil {
			return terror.NewValidationError("cannot parse cargo request uuid", "cargoRequestId")
		}
		cargoRequest, err := s.cargoRequestRepo.GetCargoRequestById(ctx, cargoRequestUUID)
		if err != nil {
			return err
		}
		cargoRequestRoutes = append(cargoRequestRoutes, cargoRequest)
	}

	cargoRequestRouteIDs := make([]string, 0, len(cargoRequestRoutes))
	for _, cargoRequestRoute := range cargoRequestRoutes {
		cargoRequestRouteIDs = append(cargoRequestRouteIDs, cargoRequestRoute.RouteID.String())
	}

	tripRouteId, err := s.client.MergeRoutes(tripId, cargoRequestRouteIDs)
	if err != nil {
		return err
	}

	s.tripRepo.UpdateTripStatus(ctx, tripId, models.TripStatusInProgress)
	s.tripRepo.UpdateRout(ctx, tripId, tripRouteId.String())
	for _, cargoRequest := range cargoRequestRoutes {
		err = s.cargoRequestRepo.UpdateCargoRequestOnStartTrip(
			ctx,
			uuid.MustParse(cargoRequest.ID.String()),
			uuid.MustParse(tripId),
			tripRouteId,
			models.CargoRequestStatusInProgress,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TripService) GetTripByIdAndCarrier(
	ctx context.Context,
	tripID *uuid.UUID,
	carrierID *int32,
) (*models.TripResponse, error) {

	trip, err := s.tripRepo.GetTripByIdAndCarrier(ctx, tripID, carrierID)
	if err != nil {
		return nil, err
	}
	if trip == nil {
		return nil, terror.NewNotFoundError("Trip", "")
	}

	fromID := pgtype.UUID{
		Bytes: trip.FromStation.Bytes,
		Valid: true,
	}
	toID := pgtype.UUID{
		Bytes: trip.ToStation.Bytes,
		Valid: true,
	}

	stations, err := s.stationRepo.GetStations(ctx, []pgtype.UUID{fromID, toID})
	if err != nil {
		return nil, err
	}

	var fromStation, toStation models.Station

	for _, st := range stations {
		id := uuid.UUID(st.ID.Bytes)

		if id == trip.FromStation.Bytes {
			fromStation = models.Station{
				Address: st.Address.String,
				Coords: models.Coords{
					Lat: st.Lat.Float64,
					Lon: st.Lon.Float64,
				},
			}
		}

		if id == trip.ToStation.Bytes {
			toStation = models.Station{
				Address: st.Address.String,
				Coords: models.Coords{
					Lat: st.Lat.Float64,
					Lon: st.Lon.Float64,
				},
			}
		}
	}

	return &models.TripResponse{
		ID:              trip.ID.Bytes,
		RouteID:         util.PgUuidToGoUuidPointer(trip.RouteID),
		FromStation:     fromStation,
		ToStation:       toStation,
		StartedAt:       trip.StartedAt.Int64,
		CalculatedEndAt: trip.CalculateEndAt.Int64,
		ActualEndAt:     trip.ActualEndAt.Int64,
		Price:           util.NumericToString(trip.Price),
		Status:          trip.Status.String,
		CarrierID:       int(trip.Carrier.Int32),
		CarID:           int(trip.Car.Int32),
	}, nil
}
