package repositories

import (
	"context"
	"fmt"
	"gruzowiki/db/pg"
	"gruzowiki/rest/models"
	"gruzowiki/util"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CargoRequestRepo struct {
	conn pg.Conn
}

func NewCargoRequestRepo(conn pg.Conn) *CargoRequestRepo {
	return &CargoRequestRepo{
		conn: conn,
	}
}

func (c *CargoRequestRepo) GetCargoRequestWithFilters(
	ctx context.Context,
	request models.GetCargoRequest,
	pageNumber int,
	pageSize int,
) ([]pg.CargoRequest, error) {
	whereSQL, args := buildCargoRequestConditions(request)

	query := "select * from cargo_requests" + whereSQL + " order by created_at desc offset $%d limit $%d"

	offset := (pageNumber - 1) * pageSize
	args = append(args, offset, pageSize)

	query = fmt.Sprintf(query, len(args)-1, len(args))

	rows, err := c.conn.Queries(ctx).RawQuery(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pg.CargoRequest
	for rows.Next() {
		var i pg.CargoRequest
		if err := rows.Scan(
			&i.ID,
			&i.ConsignerID,
			&i.RecipientID,
			&i.FromStation,
			&i.ToStation,
			&i.CreatedAt,
			&i.Deadline,
			&i.RouteID,
			&i.TripID,
			&i.Price,
			&i.Status,
			&i.ReceiveCode,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func buildCargoRequestConditions(req models.GetCargoRequest) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	if req.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *req.ID)
		argIndex++
	}

	if req.ConsignerID != nil {
		conditions = append(conditions, fmt.Sprintf("consigner_id = $%d", argIndex))
		args = append(args, *req.ConsignerID)
		argIndex++
	}

	if req.RecipientID != nil {
		conditions = append(conditions, fmt.Sprintf("recipient_id = $%d", argIndex))
		args = append(args, *req.RecipientID)
		argIndex++
	}

	if req.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if req.CreatedFrom != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, util.ToTimestamp(*req.CreatedFrom))
		argIndex++
	}

	if req.CreatedTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, util.ToTimestamp(*req.CreatedTo))
		argIndex++
	}

	if len(conditions) == 0 {
		return "", nil
	}

	whereClause := " where " + strings.Join(conditions, " and ")

	return whereClause, args
}

func (c *CargoRequestRepo) CreateCargoRequest(ctx context.Context, params pg.InsertCargoRequestParams) (*pgtype.UUID, error) {
	pgid, err := c.conn.Queries(ctx).InsertCargoRequest(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return &pgid, nil
}

func (c *CargoRequestRepo) GetCargoTypes(ctx context.Context) ([]pg.CargoType, error) {
	return c.conn.Queries(ctx).GetCargoTypes(ctx)
}

func (c *CargoRequestRepo) CreateCargo(ctx context.Context, cargos []models.Cargo) ([]pgtype.UUID, error) {
	ids := make([]pgtype.UUID, len(cargos))
	for i, cargo := range cargos {
		uuid, err := uuid.Parse(cargo.CargoRequestId)
		if err != nil {
			return nil, err
		}
		id, err := c.conn.Queries(ctx).CreateCargo(ctx, pg.CreateCargoParams{
			Length:    pgtype.Int4{Int32: int32(cargo.Length), Valid: true},
			Width:     pgtype.Int4{Int32: int32(cargo.Width), Valid: true},
			Height:    pgtype.Int4{Int32: int32(cargo.Height), Valid: true},
			Weight:    pgtype.Int4{Int32: int32(cargo.Weight), Valid: true},
			CargoType: pgtype.Int4{Int32: int32(cargo.CargoType), Valid: true},
			Worth:     pgtype.Numeric{Int: big.NewInt(int64(cargo.Worth)), Valid: true},
			RequestID: pgtype.UUID{Bytes: uuid, Valid: true},
		})
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}

	return ids, nil
}

func (c *CargoRequestRepo) MarkTrip(ctx context.Context, cargoReqId string, tripId string) error {
	id, err := uuid.Parse(cargoReqId)
	if err != nil {
		return err
	}
	trip, err := uuid.Parse(tripId)
	if err != nil {
		return err
	}
	_, err = c.conn.Queries(ctx).UpdateCargoRequestTrip(ctx, pg.UpdateCargoRequestTripParams{
		ID:     pgtype.UUID{Bytes: id, Valid: true},
		TripID: pgtype.UUID{Bytes: trip, Valid: true},
	})
	return err
}

func (c *CargoRequestRepo) UpdateCargoRequestCode(ctx context.Context, reqId, code string) error {
	/*id, err := uuid.Parse(reqId)
	if err != nil {
		return err
	}
	receiveCode, err := strconv.ParseInt(code, 10, 32)
	if err != nil {
		return err
	}
	err = c.conn.Queries(ctx).UpdateCargoRequestReceiveCode(ctx, pg.UpdateCargoRequestReceiveCodeParams{
		ID:          pgtype.UUID{Bytes: id, Valid: true},
		ReceiveCode: pgtype.Int4{Int32: int32(receiveCode), Valid: true},
	})*/
	return nil
}

func (c *CargoRequestRepo) UpdateRoute(ctx context.Context, cargoReqId string, routeId string) error {
	id, err := uuid.Parse(cargoReqId)
	if err != nil {
		return err
	}
	route, err := uuid.Parse(routeId)
	if err != nil {
		return err
	}
	err = c.conn.Queries(ctx).UpdateCargoRequestRoute(ctx, pg.UpdateCargoRequestRouteParams{
		ID:      pgtype.UUID{Bytes: id, Valid: true},
		RouteID: pgtype.UUID{Bytes: route, Valid: true},
	})
	return err
}

type CargoRequestPair struct {
	CargoRequestID string
	RouteID        string
}

func (c *CargoRequestRepo) GetCargoRequestIDAndRoute(ctx context.Context, cargoReqId string) (CargoRequestPair, error) {
	id, err := uuid.Parse(cargoReqId)
	if err != nil {
		return CargoRequestPair{}, err
	}
	req, err := c.conn.Queries(ctx).GetCargoRequestIDAndRoute(ctx, pgtype.UUID{Bytes: id, Valid: true})
	return CargoRequestPair{req.ID.String(), req.RouteID.String()}, err
}

func (c *CargoRequestRepo) SetRouteID(ctx context.Context, cargoRequestId uuid.UUID, routeID uuid.UUID) error {
	return c.conn.Queries(ctx).SetRouteIDForCargoRequest(ctx, pg.SetRouteIDForCargoRequestParams{
		ID:      pgtype.UUID{Bytes: cargoRequestId, Valid: true},
		RouteID: pgtype.UUID{Bytes: routeID, Valid: true},
	})
}

func (c *CargoRequestRepo) SetTrip(ctx context.Context, cargoRequestId uuid.UUID, tripId uuid.UUID) error {
	return c.conn.Queries(ctx).SetTripIDForCargoRequest(ctx, pg.SetTripIDForCargoRequestParams{
		ID:     pgtype.UUID{Bytes: cargoRequestId, Valid: true},
		TripID: pgtype.UUID{Bytes: tripId, Valid: true},
	})
}

func (c *CargoRequestRepo) UpdateCargoRequestOnStartTrip(ctx context.Context, cargoRequestId uuid.UUID, tripId uuid.UUID, routeId uuid.UUID, status string) error {
	return c.conn.Queries(ctx).UpdateCargoRequestOnStartTrip(ctx, pg.UpdateCargoRequestOnStartTripParams{
		ID:      pgtype.UUID{Bytes: cargoRequestId, Valid: true},
		TripID:  pgtype.UUID{Bytes: tripId, Valid: true},
		RouteID: pgtype.UUID{Bytes: routeId, Valid: true},
		Status:  pgtype.Text{String: status, Valid: true},
	})
}

func (c *CargoRequestRepo) GetCargoRequestsForTrip(
	ctx context.Context,
	maxLength *int32,
	maxWidth *int32,
	maxHeight *int32,
	cargoType *int32,
	deadline *int64,
	minPrice *pgtype.Numeric,
) ([]pgtype.UUID, error) {

	arg := pg.GetCargoRequestsForTripParams{}

	if maxLength != nil {
		arg.MaxLength = pgtype.Int4{Int32: *maxLength, Valid: true}
	} else {
		arg.MaxLength = pgtype.Int4{Valid: false}
	}

	if maxWidth != nil {
		arg.MaxWidth = pgtype.Int4{Int32: *maxWidth, Valid: true}
	} else {
		arg.MaxWidth = pgtype.Int4{Valid: false}
	}

	if maxHeight != nil {
		arg.MaxHeight = pgtype.Int4{Int32: *maxHeight, Valid: true}
	} else {
		arg.MaxHeight = pgtype.Int4{Valid: false}
	}

	if cargoType != nil {
		arg.CargoType = pgtype.Int4{Int32: *cargoType, Valid: true}
	} else {
		arg.CargoType = pgtype.Int4{Valid: false}
	}

	if deadline != nil {
		arg.Deadline = pgtype.Int8{Int64: *deadline, Valid: true}
	} else {
		arg.Deadline = pgtype.Int8{Valid: false}
	}

	if minPrice != nil {
		arg.MinPrice = *minPrice
	} else {
		arg.MinPrice = pgtype.Numeric{Valid: false}
	}

	return c.conn.Queries(ctx).GetCargoRequestsForTrip(ctx, arg)
}

func (c *CargoRequestRepo) GetRequestsRouteIds(
	ctx context.Context,
	ids []pgtype.UUID,
) ([]pg.GetCargoRequestRouteIDsRow, error) {
	return c.conn.Queries(ctx).GetCargoRequestRouteIDs(ctx, ids)
}

func (c *CargoRequestRepo) GetTripRouteId(
	ctx context.Context,
	tripID pgtype.UUID,
) (pgtype.UUID, error) {
	return c.conn.Queries(ctx).GetTripRouteID(ctx, tripID)
}

func (c *CargoRequestRepo) GetCargoRequestById(ctx context.Context, id uuid.UUID) (pg.CargoRequest, error) {
	pgID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	req, err := c.conn.Queries(ctx).GetCargoRequest(ctx, pgID)
	if err != nil {
		return pg.CargoRequest{}, err
	}

	return req, nil
}

func (c *CargoRequestRepo) GetCargoRequestByIdWithCargo(
	ctx context.Context,
	id uuid.UUID,
) (pg.GetCargoRequestWithCargoRow, error) {
	return c.conn.Queries(ctx).GetCargoRequestWithCargo(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
}
