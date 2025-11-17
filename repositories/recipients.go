package repositories

import (
	"context"
	"errors"
	"fmt"
	"gruzowiki/db/pg"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RecipientRepo struct {
	conn pg.Conn
}

func NewRecipientRepo(conn pg.Conn) *RecipientRepo {
	return &RecipientRepo{conn: conn}
}

func (r *RecipientRepo) CreateRecipient(ctx context.Context, firstName, secondName, thirdName, phone, email string) (int32, error) {
	params := pg.CreateRecipientParams{
		FirstName:  pgtype.Text{String: firstName, Valid: true},
		SecondName: pgtype.Text{String: secondName, Valid: true},
		ThirdName:  pgtype.Text{String: thirdName, Valid: true},
		Phone:      pgtype.Text{String: phone, Valid: true},
		Email:      pgtype.Text{String: email, Valid: true},
	}
	
	newID, err := r.conn.Queries(ctx).CreateRecipient(ctx, params)
	if err != nil {
		fmt.Printf("[Repo] Error from DB: %v\n", err)
		return 0, fmt.Errorf("query CreateRecipient: %w", err)
	}

	return newID, nil
}

func (r *RecipientRepo) GetRecipientByID(ctx context.Context, id int32) (*pg.Recipient, error) {
	recipient, err := r.conn.Queries(ctx).GetRecipient(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query GetRecipient: %w", err)
	}
	return &recipient, nil
}

func (r *RecipientRepo) ListRecipients(ctx context.Context) ([]pg.Recipient, error) {
	items, err := r.conn.Queries(ctx).ListRecipients(ctx)
	if err != nil {
		return nil, fmt.Errorf("query ListRecipients: %w", err)
	}
	return items, nil
}

func (r *RecipientRepo) UpdateRecipient(ctx context.Context, id int32, firstName, secondName, thirdName, phone, email string) (int32, error) {
	params := pg.UpdateRecipientParams{
		ID:         id,
		FirstName:  pgtype.Text{String: firstName, Valid: true},
		SecondName: pgtype.Text{String: secondName, Valid: true},
		ThirdName:  pgtype.Text{String: thirdName, Valid: true},
		Phone:      pgtype.Text{String: phone, Valid: true},
		Email:      pgtype.Text{String: email, Valid: true},
	}

	updatedID, err := r.conn.Queries(ctx).UpdateRecipient(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("query UpdateRecipient: %w", err)
	}
	return updatedID, nil
}

func (r *RecipientRepo) DeleteRecipient(ctx context.Context, id int32) error {
    err := r.conn.Queries(ctx).DeleteRecipient(ctx, id)
    if err != nil {
        return fmt.Errorf("query DeleteRecipient: %w", err)
    }
    return nil
}