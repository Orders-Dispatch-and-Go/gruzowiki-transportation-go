package util

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UuidToPgUuid(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func PgUuidToUuid(id pgtype.UUID) uuid.UUID {
	return id.Bytes
}

func GoTextToPgText(text string) pgtype.Text {
	return pgtype.Text{String: text, Valid: true}
}
