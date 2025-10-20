package carrier

import (
	"auth-service/internal/api/controller"
	"auth-service/internal/domain/carrierDto"
	"context"
)

type CarrierController interface {
	GetCarrier(
		id int32,
		ctx context.Context,
		requestReader controller.RequestReader,
	) (getCarrierResponse carrierDto.GetCarrierResponse, err error)
}
