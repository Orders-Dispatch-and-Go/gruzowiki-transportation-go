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
	//if !n.Valid {
	//	return ""
	//}
	//
	//if n.Exp == 0 {
	//	return n.Int.String()
	//}
	//
	//rat := new(big.Rat).SetInt(n.Int)
	//
	//scale := new(big.Rat).SetFloat64(math.Pow10(int(-n.Exp)))
	//rat = rat.Mul(rat, scale)

	//return rat.FloatString(int(-n.Exp))
	numLen := len(n.Int.String())
	return n.Int.String()[:numLen-2] + "." + n.Int.String()[numLen-2:]
}

func PgInt8ToInt(n pgtype.Int8) int64 {
	return n.Int64
}

func PgInt4ToInt(n pgtype.Int4) int {
	return int(n.Int32)
}
