package services

import (
	"context"
	"gruzowiki/db/pg"
	"gruzowiki/rest/middlewares"
	"gruzowiki/rest/models"
	"gruzowiki/util"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	CargoRequestService struct {
		repo           CargoRequestRepo
		stationService StationService
	}

	CargoRequestRepo interface {
		CreateCargoRequest(ctx context.Context, params pg.InsertCargoRequestParams) (*pgtype.UUID, error)
		GetCargoRequestWithFilters(
			ctx context.Context,
			request models.GetCargoRequest,
			pageNumber int,
			pageSize int,
		) ([]pg.CargoRequest, error)
		GetCargoTypes(ctx context.Context) ([]pg.CargoType, error)
		CreateCargo(ctx context.Context, cargos []models.Cargo) ([]pgtype.UUID, error)
		MarkTrip(ctx context.Context, cargoReqId string, tripId string) error
	}

	StationService interface {
		GetStations(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]models.Station, error)
		CreateStation(ctx context.Context, station models.Station) (*models.CreateStationResponse, error)
	}
)

func NewCargoRequestService(repo CargoRequestRepo, stationService StationService) *CargoRequestService {
	return &CargoRequestService{
		repo:           repo,
		stationService: stationService,
	}
}

func (s *CargoRequestService) SearchCargoRequests(
	ctx context.Context,
	filter models.GetCargoRequest,
	pageNumber int,
	pageSize int,
) (*models.SearchCargoRequestsResponse, error) {
	pgCargoRequests, err := s.repo.GetCargoRequestWithFilters(ctx, filter, pageNumber, pageSize)

	if err != nil {
		return nil, err
	}

	requestUserId := ctx.Value(middlewares.UserIdCtxClaim).(int)

	responseCargoRequests := make([]models.CargoRequestResponse, 0, len(pgCargoRequests))
	for _, pgReq := range pgCargoRequests {
		fromStationId := util.PgUuidToUuid(pgReq.FromStation)
		toStationId := util.PgUuidToUuid(pgReq.ToStation)
		stations, err := s.stationService.GetStations(ctx, []uuid.UUID{fromStationId, toStationId})

		if err != nil {
			return nil, err
		}

		fromStation := stations[fromStationId]
		toStation := stations[toStationId]

		id, _ := uuid.FromBytes(pgReq.ID.Bytes[:])
		maxPrice := util.NumericToString(pgReq.Price)

		var receiveCode *string = nil
		if requestUserId == util.PgInt4ToInt(pgReq.ConsignerID) {
			stringReceiveCode := strconv.Itoa(util.PgInt4ToInt(pgReq.ReceiveCode))
			receiveCode = &stringReceiveCode
		}

		responseCargoRequests = append(responseCargoRequests, models.CargoRequestResponse{
			ID:          id.String(),
			ConsignerID: util.PgInt4ToInt(pgReq.ConsignerID),
			RecipientID: util.PgInt4ToInt(pgReq.RecipientID),
			FromStation: &fromStation,
			ToStation:   &toStation,
			CreatedAt:   util.PgInt8ToInt(pgReq.CreatedAt),
			Deadline:    util.PgInt8ToInt(pgReq.Deadline),
			RouteID:     nil,
			TripID:      nil,
			Price:       maxPrice,
			Status:      pgReq.Status.String,
			ReceiveCode: receiveCode,
		})
	}

	return &models.SearchCargoRequestsResponse{
		CargoRequests: responseCargoRequests,
	}, nil
}

func (s *CargoRequestService) CreateCargoRequest(ctx context.Context, postCargoRequestRequest models.PostCargoRequestRequest) (*models.PostCargoRequestResponse, error) {
	createFromStationResponse, err := s.stationService.CreateStation(ctx, postCargoRequestRequest.FromStation)
	if err != nil {
		return nil, err
	}

	createToStationResponse, err := s.stationService.CreateStation(ctx, postCargoRequestRequest.ToStation)
	if err != nil {
		return nil, err
	}

	receiveCode, err := util.GenerateRandomReceiveCode()
	if err != nil {
		return nil, err
	}

	id, err := s.repo.CreateCargoRequest(ctx, pg.InsertCargoRequestParams{
		ID:          util.UuidToPgUuid(uuid.New()),
		ConsignerID: util.Int32ToPgInt4(postCargoRequestRequest.ConsignerID),
		RecipientID: util.Int32ToPgInt4(postCargoRequestRequest.RecipientID),
		FromStation: util.UuidToPgUuid(createFromStationResponse.ID),
		ToStation:   util.UuidToPgUuid(createToStationResponse.ID),
		CreatedAt:   util.Int64ToPgInt8(time.Now().Unix()),
		Deadline:    util.Int64ToPgInt8(util.ToTimestamp(postCargoRequestRequest.Deadline)),
		Price:       util.ToNumeric(postCargoRequestRequest.MaxPrice),
		Status:      util.GoTextToPgText(models.CargoRequestStatusPending),
		ReceiveCode: util.Int32ToPgInt4(receiveCode),
	})

	if err != nil {
		return nil, err
	}

	return &models.PostCargoRequestResponse{ID: id.Bytes, ReceiveCode: strconv.Itoa(int(receiveCode))}, nil
}

func (s *CargoRequestService) GetCargoTypes(ctx context.Context) ([]models.CargoType, error) {
	types, err := s.repo.GetCargoTypes(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]models.CargoType, len(types))
	for i, t := range types {
		resp[i] = models.CargoType{
			Id:      int(t.ID),
			Type:    t.Type.String,
			Fragile: t.Fragile.Bool,
		}
	}
	return resp, nil
}

func (s *CargoRequestService) CreateCargo(ctx context.Context, cargo []models.Cargo) ([]string, error) {
	ids, err := s.repo.CreateCargo(ctx, cargo)
	if err != nil {
		return nil, err
	}

	resp := make([]string, 0, len(ids))
	for i, _ := range ids {
		resp[i] = uuid.UUID(ids[i].Bytes).String()
	}

	return resp, nil
}

func (s *CargoRequestService) MarkTrip(ctx context.Context, cargoReqId string, tripId string) error {
	return s.repo.MarkTrip(ctx, cargoReqId, tripId)
}
