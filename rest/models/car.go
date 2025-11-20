package models

type CreateCarRequest struct {
	Type      string `json:"type"`
	Length    int32  `json:"length"`
	Width     int32  `json:"width"`
	Height    int32  `json:"height"`
	MaxWeight int32  `json:"maxweight"`
	Number    string `json:"number"`
	OwnerID   int32  `json:"ownerid"`
}

type UpdateCarRequest struct {
	Type      *string `json:"type"`
	Length    *int32  `json:"length"`
	Width     *int32  `json:"width"`
	Height    *int32  `json:"height"`
	MaxWeight *int32  `json:"maxweight"`
	Number    *string `json:"number"`
	OwnerID   *int32  `json:"ownerid"`
}

type CreateCarResponse struct {
	ID int32 `json:"id"`
}

type GetCarResponse struct {
	ID        int32  `json:"id"`
	Type      string `json:"type"`
	Length    int32  `json:"length"`
	Width     int32  `json:"width"`
	Height    int32  `json:"height"`
	MaxWeight int32  `json:"maxweight"`
	Number    string `json:"number"`
	OwnerID   int32  `json:"ownerid"`
}

type UpdateCarResponse struct {
	ID int32 `json:"id"`
}