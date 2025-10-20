package carrier

import (
	"auth-service/internal/api/controller"
	"auth-service/internal/api/controller/carrier"
	"context"
)

type CarrierService interface {
	GetCarrier(id int32, ctx context.Context, requestReader controller.RequestReader) (
		getCarrierResponse carrier.GetCarrierResponse,
		err error,
	)
}
