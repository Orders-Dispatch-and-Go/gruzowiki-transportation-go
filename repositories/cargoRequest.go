package repositories

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"gruzowiki/db/pg"
	"gruzowiki/rest/models"
	"gruzowiki/util"
	"math/big"
	"strings"
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
