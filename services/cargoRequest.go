package services

import (
	"context"
	"errors"
	"fmt"
	"gruzowiki/db/pg"
	"gruzowiki/rest/middlewares"
	"gruzowiki/rest/models"
	"gruzowiki/util"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	CargoRequestService struct {
		repo           CargoRequestRepo
		stationService StationService
		client         CargoRequestFeignClient
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
		UpdateCargoRequestCode(ctx context.Context, cargoReqId string, code string) error
		GetCargoRequestsForTrip(
			ctx context.Context,
			maxLength *int32,
			maxWidth *int32,
			maxHeight *int32,
			cargoType *int32,
			deadline *int64,
			minPrice *pgtype.Numeric,
		) ([]pgtype.UUID, error)
		GetTripRouteId(ctx context.Context, tripID pgtype.UUID) (pgtype.UUID, error)
		GetRequestsRouteIds(
			ctx context.Context,
			ids []pgtype.UUID,
		) ([]pg.GetCargoRequestRouteIDsRow, error)
		GetCargoRequestById(ctx context.Context, id uuid.UUID) (pg.CargoRequest, error)
	}

	StationService interface {
		GetStations(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]models.Station, error)
		CreateStation(ctx context.Context, station models.Station) (*models.CreateStationResponse, error)
	}

	CargoRequestFeignClient interface {
		CreateRouteForCargoRequest(request models.PostCargoRequestRequest, fromStationId uuid.UUID, toStationId uuid.UUID) (*uuid.UUID, error)
		GetPotentialTrips(tripRouteID string, cargoRequestRouteIDs []string) ([]string, error)
	}
)

func NewCargoRequestService(repo CargoRequestRepo, stationService StationService, client CargoRequestFeignClient) *CargoRequestService {
	return &CargoRequestService{
		repo:           repo,
		stationService: stationService,
		client:         client,
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

		routeId := util.PgUuidToUuid(pgReq.RouteID).String()
		tripId := util.PgUuidToUuid(pgReq.TripID).String()

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
			RouteID:     &routeId,
			TripID:      &tripId,
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

	routeId, err := s.client.CreateRouteForCargoRequest(
		postCargoRequestRequest,
		createFromStationResponse.ID,
		createToStationResponse.ID,
	)
	if err != nil {
		return nil, err
	}

	requestUserId := ctx.Value(middlewares.UserIdCtxClaim).(int)

	id, err := s.repo.CreateCargoRequest(ctx, pg.InsertCargoRequestParams{
		ID:          util.UuidToPgUuid(uuid.New()),
		ConsignerID: util.Int32ToPgInt4(int32(requestUserId)),
		RecipientID: util.Int32ToPgInt4(postCargoRequestRequest.RecipientID),
		FromStation: util.UuidToPgUuid(createFromStationResponse.ID),
		ToStation:   util.UuidToPgUuid(createToStationResponse.ID),
		CreatedAt:   util.Int64ToPgInt8(time.Now().Unix()),
		Deadline:    util.Int64ToPgInt8(util.ToTimestamp(postCargoRequestRequest.Deadline)),
		RouteID:     util.UuidToPgUuid(*routeId),
		Price:       util.ToNumeric(postCargoRequestRequest.MaxPrice),
		Status:      util.GoTextToPgText(models.CargoRequestStatusPending),
		ReceiveCode: util.Int32ToPgInt4(receiveCode),
	})

	if err != nil {
		return nil, err
	}

	return &models.PostCargoRequestResponse{ID: id.Bytes, ReceiveCode: strconv.Itoa(int(receiveCode))}, nil
}

func (s *CargoRequestService) GetRequestsForTripWithRoutes(
	ctx context.Context,
	tripID uuid.UUID,
	filter models.GetCargoRequestsForTripFilter,
) (*models.GetCargoRequestsForTripResponse, error) {
	var minPrice *pgtype.Numeric
	if filter.MinPrice != nil {
		n := util.ToNumeric(decimal.NewFromInt(*filter.MinPrice))
		minPrice = &n
	}

	//получили тут id заявок, который нам подходят по фильтрам
	ids, err := s.repo.GetCargoRequestsForTrip(
		ctx,
		filter.CargoLengthMax,
		filter.CargoWidthMax,
		filter.CargoHeightMax,
		filter.CargoType,
		filter.Deadline,
		minPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("get cargo requests by params: %w", err)
	}

	// если ничего не найдено - выходим
	if len(ids) == 0 {
		return &models.GetCargoRequestsForTripResponse{CargoRequests: []models.CargoRequestResponse{}}, nil
	}

	//получаем по tripId из query - id маршрута для передачи в rust potential
	pgTripID := util.UuidToPgUuid(tripID)
	tripRoutePg, err := s.repo.GetTripRouteId(ctx, pgTripID)
	if err != nil {
		return nil, fmt.Errorf("get trip route id: %w", err)
	}

	if !tripRoutePg.Valid {
		return nil, errors.New("trip has no route_id")
	}

	tRidUuid, err := uuid.FromBytes(tripRoutePg.Bytes[:])
	if err != nil {
		return nil, fmt.Errorf("invalid trip route id bytes: %w", err)
	}

	//получаем id маршрутов для отфильтрованных заявок для передачи в rust potential
	reqRoutes, err := s.repo.GetRequestsRouteIds(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get cargo request route ids: %w", err)
	}
	if len(reqRoutes) == 0 {
		return &models.GetCargoRequestsForTripResponse{CargoRequests: []models.CargoRequestResponse{}}, nil
	}

	//делаем мапу routeToReqs для соответствия id заявки - id маршрута. Т.к после получения id маршрутов из potential
	//надо убрать ненужные маршруты и соотвествующие им отфильтрованные заявки
	routeIDs := make([]string, 0, len(reqRoutes))
	routeToReqs := make(map[string][]uuid.UUID, len(reqRoutes))

	for _, row := range reqRoutes {
		if !row.RouteID.Valid || row.RouteID.Bytes == [16]byte{} {
			continue
		}
		rID, err := uuid.FromBytes(row.RouteID.Bytes[:])
		if err != nil {
			continue
		}
		routeIDs = append(routeIDs, rID.String())

		if row.ID.Valid {
			reqID, err := uuid.FromBytes(row.ID.Bytes[:])
			if err == nil {
				routeToReqs[rID.String()] = append(routeToReqs[rID.String()], reqID)
			}
		}
	}

	//идем к rust potential получать нужные нам заявки по географии
	matchingRouteIDs, err := s.client.GetPotentialTrips(tRidUuid.String(), routeIDs)
	if err != nil {
		return nil, fmt.Errorf("routing error: %w", err)
	}
	if len(matchingRouteIDs) == 0 {
		return &models.GetCargoRequestsForTripResponse{CargoRequests: []models.CargoRequestResponse{}}, nil
	}

	routeSet := make(map[string]struct{}, len(matchingRouteIDs))
	for _, rid := range matchingRouteIDs {
		routeSet[rid] = struct{}{}
	}

	requestsToLoad := make([]uuid.UUID, 0)
	for routeID, reqs := range routeToReqs {
		if _, ok := routeSet[routeID]; !ok {
			continue
		}
		for _, rid := range reqs {
			requestsToLoad = append(requestsToLoad, rid)
		}
	}

	if len(requestsToLoad) == 0 {
		return &models.GetCargoRequestsForTripResponse{CargoRequests: []models.CargoRequestResponse{}}, nil
	}

	finalRequests := make([]models.CargoRequestResponse, 0, len(requestsToLoad))
	for _, reqID := range requestsToLoad { //подгружаем нужные нам заявки по id
		pgReq, err := s.repo.GetCargoRequestById(ctx, reqID)
		if err != nil {
			return nil, fmt.Errorf("load cargo request %s: %w", reqID.String(), err)
		}

		var routeIDStr *string
		if pgReq.RouteID.Valid {
			r := util.PgUuidToUuid(pgReq.RouteID).String()
			routeIDStr = &r
		}

		var tripIDStr *string
		if pgReq.TripID.Valid {
			t := util.PgUuidToUuid(pgReq.TripID).String()
			tripIDStr = &t
		}

		var receiveCode *string
		if pgReq.ReceiveCode.Valid {
			rc := strconv.Itoa(int(pgReq.ReceiveCode.Int32))
			receiveCode = &rc
		}

		fromStationID := util.PgUuidToUuid(pgReq.FromStation)
		toStationID := util.PgUuidToUuid(pgReq.ToStation)

		stations, err := s.stationService.GetStations(ctx, []uuid.UUID{fromStationID, toStationID})
		if err != nil {
			return nil, err
		}

		fromStation := stations[fromStationID]
		toStation := stations[toStationID]

		req := models.CargoRequestResponse{
			ID:          reqID.String(),
			ConsignerID: int(pgReq.ConsignerID.Int32),
			RecipientID: int(pgReq.RecipientID.Int32),
			FromStation: &fromStation,
			ToStation:   &toStation,
			CreatedAt:   pgReq.CreatedAt.Int64,
			Deadline:    pgReq.Deadline.Int64,
			RouteID:     routeIDStr,
			TripID:      tripIDStr,
			Price:       util.NumericToString(pgReq.Price),
			Status:      pgReq.Status.String,
			ReceiveCode: receiveCode,
		}

		finalRequests = append(finalRequests, req)
	}

	return &models.GetCargoRequestsForTripResponse{
		CargoRequests: finalRequests,
	}, nil
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

	resp := make([]string, len(ids))
	for i, _ := range ids {
		resp[i] = uuid.UUID(ids[i].Bytes).String()
	}

	return resp, nil
}

func (s *CargoRequestService) MarkTrip(ctx context.Context, cargoReqId string, tripId string) error {
	return s.repo.MarkTrip(ctx, cargoReqId, tripId)
}

func (s *CargoRequestService) Delivered(ctx context.Context, cargoReqId string, code string) error {
	return s.repo.UpdateCargoRequestCode(ctx, cargoReqId, code)
}
