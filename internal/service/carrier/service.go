package carrier

import (
	"auth-service/internal/api/controller/carrier"
	"context"
)

type Service struct {
	// соединение с бд
	// что то еще 100%
}

func GetCarrier(id int32, ctx context.Context) (
	getCarrierResponse carrier.GetCarrierResponse,
	err error,
) {
	return carrier.GetCarrierResponse{}, nil
}
