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
