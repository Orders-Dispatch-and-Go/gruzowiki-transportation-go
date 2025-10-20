package carrier

import (
	"auth-service/internal/api/controller"
	"auth-service/internal/domain/carrierDto"
	"auth-service/internal/service/carrier"
	"context"
)

type Controller struct {
	requestReader controller.RequestReader
	service       carrier.CarrierService
}

func New(requestReader controller.RequestReader, service carrier.CarrierService) *Controller {
	return &Controller{requestReader, service}
}

func (c Controller) GetCarrier(
	id int32,
	ctx context.Context,
	requestReader controller.RequestReader,
) (getCarrierResponse carrierDto.GetCarrierResponse, err error) {
	//TODO implement me
	panic("implement me")
}
