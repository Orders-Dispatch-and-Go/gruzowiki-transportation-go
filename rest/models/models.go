package models

type CreateCarRequest struct {
	Type      string
	Length    int32
	Width     int32
	Height    int32
	MaxWeight int32
	Number    string
	OwnerID   int32
}

type CreateCarResponse struct {
	ID int32
}

type GetCarResponse struct {
	ID        int32
	Type      string
	Length    int32
	Width     int32
	Height    int32
	MaxWeight int32
	Number    string
	OwnerID   int32
}

type UpdateCarRequest struct {
	Type      *string
	Length    *int32
	Width     *int32
	Height    *int32
	MaxWeight *int32
	Number    *string
	OwnerID   *int32
}

type UpdateCarResponse struct {
	ID int32
}

type CreateCarrierRequest struct {
	ID             int32
	DriverCategory string
}

type CreateCarrierResponse struct {
	ID int32
}

type GetCarrierResponse struct {
	ID             int32
	DriverCategory string
}

type UpdateCarrierRequest struct {
	DriverCategory string
}

type UpdateCarrierResponse struct {
	ID int32
}
