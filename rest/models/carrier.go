package models

type CreateCarrierRequest struct {
	ID             int32  `json:"id"`
	DriverCategory string `json:"drivercategory"`
}

type UpdateCarrierRequest struct {
	DriverCategory string `json:"drivercategory"`
}

type CreateCarrierResponse struct {
	ID int32 `json:"id"`
}

type GetCarrierResponse struct {
	ID             int32  `json:"id"`
	DriverCategory string `json:"drivercategory"`
}

type UpdateCarrierResponse struct {
	ID int32 `json:"id"`
}
