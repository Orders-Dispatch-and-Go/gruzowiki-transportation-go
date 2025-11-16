package util

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func ToNumeric(decimalNum decimal.Decimal) pgtype.Numeric {
	price := pgtype.Numeric{}
	_ = price.Scan(decimalNum)
	return price
}
