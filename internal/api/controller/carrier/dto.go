package carrier

import "github.com/jackc/pgx/v5/pgtype"

type GetCarrierResponse struct {
	ID             int32
	DriverCategory pgtype.Text
}
