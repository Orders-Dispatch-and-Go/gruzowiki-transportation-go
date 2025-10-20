package carrierRepo

import (
	"auth-service/internal/db/pg"
	"context"
)

type CarrierRepo interface {
	GetCarrierById(id int32, ctx context.Context) (*pg.Carrier, error)
}
