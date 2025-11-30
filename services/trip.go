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

type RouterserviceClient interface {
    GetPotentialTrips(cargoRequestRouteID string, tripRouteIDs []string) ([]string, error)
}

type TripService struct {
    tripRepo    *repositories.TripRepo
    stationRepo *repositories.StationRepo
    routerserviceClient RouterserviceClient
}

func NewTripService(tripRepo *repositories.TripRepo, stationRepo *repositories.StationRepo, routerserviceClient RouterserviceClient) *TripService {
    return &TripService{
        tripRepo:            tripRepo,
        stationRepo:         stationRepo,
        routerserviceClient: routerserviceClient,
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

func (s *TripService) GetTripsByCargoRequest(ctx context.Context, cargoRequestID uuid.UUID, pageNumber, pageSize int) (models.Trips, error) {
    ids, err := s.tripRepo.GetTripsIdByCargoRequest(ctx, cargoRequestID) 
    if err != nil {
        return models.Trips{}, err
    }

    suitsTripsIds, err := s.routerserviceClient.GetPotentialTrips(cargoRequestID.String(), uuidsToStrigs(ids))
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
            ID:              uuid.UUID(trip.ID.Bytes),
            FromStation:     models.Station{
                Address: trip.FromAddress.String,
                Coords:  models.Coords{Lat: trip.FromLat.Float64, Lon: trip.FromLon.Float64},
            },
            ToStation:       models.Station{
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
