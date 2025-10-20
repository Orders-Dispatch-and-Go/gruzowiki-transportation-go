package carrier

import (
	"auth-service/internal/api/controller"
	"context"
)

type Controller struct {
	requestReader controller.RequestReader
	controller    CarrierController
}

func New(requestReader controller.RequestReader, controller Controller) *Controller {
	return &Controller{
		requestReader: requestReader,
		controller:    controller,
	}
}

func (c Controller) GetCarrier(
	id int32,
	ctx context.Context,
	requestReader controller.RequestReader,
) (getCarrierResponse GetCarrierResponse, err error) {
	//TODO implement me
	panic("implement me")
}
