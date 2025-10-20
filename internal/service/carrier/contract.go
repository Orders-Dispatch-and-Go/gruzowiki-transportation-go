package carrier

import (
	"auth-service/internal/domain/carrierDto"
	"context"
)

type CarrierService interface {
	GetCarrier(id int32, ctx context.Context) (
		getCarrierResponse carrierDto.GetCarrierResponse,
		err error,
	)
}
