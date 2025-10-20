package carrier

import (
	"auth-service/internal/domain/carrierDto"
	"context"
)

type Service struct {
	// соединение с бд
	// что то еще 100%
}

func GetCarrier(id int32, ctx context.Context) (
	getCarrierResponse carrierDto.GetCarrierResponse,
	err error,
) {
	return carrierDto.GetCarrierResponse{}, nil
}
