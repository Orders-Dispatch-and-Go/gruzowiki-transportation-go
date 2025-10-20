package carrier

import (
	"auth-service/internal/api/controller"
	"context"
)

type CarrierController interface {
	GetCarrier(id int32, ctx context.Context, requestReader controller.RequestReader) (getCarrierResponse GetCarrierResponse, err error)
}
