package util

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func ToNumeric(decimalNum decimal.Decimal) pgtype.Numeric {
	price := pgtype.Numeric{}
	_ = price.Scan(decimalNum.String())
	return price
}

func NumericToString(n pgtype.Numeric) string {
	numLen := len(n.Int.String())
	return n.Int.String()[:numLen-2] + "." + n.Int.String()[numLen-2:]
}

func PgInt8ToInt(n pgtype.Int8) int64 {
	return n.Int64
}

func PgInt4ToInt(n pgtype.Int4) int {
	return int(n.Int32)
}

func Float64ToPgFloat8(n float64) pgtype.Float8 {
	return pgtype.Float8{Float64: n, Valid: true}
}

func PgFloat8ToFloat64(n pgtype.Float8) float64 {
	return n.Float64
}

func Int64ToPgInt8(n int64) pgtype.Int8 {
	return pgtype.Int8{Int64: n, Valid: true}
}

func Int32ToPgInt4(n int32) pgtype.Int4 {
	return pgtype.Int4{Int32: n, Valid: true}
}

func PgInt4ToInt32(n pgtype.Int4) int32 {
	return n.Int32
}
