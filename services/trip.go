package services

import (
    "context"
    "gruzowiki/repositories"
    "gruzowiki/rest/models"
    "gruzowiki/rest/terror"
    "gruzowiki/util"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgtype"
)

type TripService struct {
    tripRepo    *repositories.TripRepo
    stationRepo *repositories.StationRepo
}

func NewTripService(tripRepo *repositories.TripRepo, stationRepo *repositories.StationRepo) *TripService {
    return &TripService{
        tripRepo:    tripRepo,
        stationRepo: stationRepo,
    }
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
        if id == uuid.UUID(trip.FromStation.Bytes) {
            fromStation = models.Station{
                Address: st.Address.String,
                Coords:  models.Coords{Lat: st.Lat.Float64, Lon: st.Lon.Float64},
            }
        }
        if id == uuid.UUID(trip.ToStation.Bytes) {
            toStation = models.Station{
                Address: st.Address.String,
                Coords:  models.Coords{Lat: st.Lat.Float64, Lon: st.Lon.Float64},
            }
        }
    }

    return &models.TripResponse{
        ID:              uuid.UUID(trip.ID.Bytes),
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
