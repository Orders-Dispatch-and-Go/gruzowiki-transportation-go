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

func PgUuidToGoUuidPointer(pguuid pgtype.UUID) *uuid.UUID {
	if pguuid.Valid == false {
		return nil
	}
	retid := uuid.UUID(pguuid.Bytes)
	return &retid
}

func PgInt4ToGoInt32Pointer(num pgtype.Int4) *int32 {
	if num.Valid == false {
		return nil
	}
	return &num.Int32
}
