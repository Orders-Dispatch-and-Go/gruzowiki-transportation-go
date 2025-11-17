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

type UpdateCarRequest struct {
	Type      *string `json:"type"`
	Length    *int32  `json:"length"`
	Width     *int32  `json:"width"`
	Height    *int32  `json:"height"`
	MaxWeight *int32  `json:"maxweight"`
	Number    *string `json:"number"`
	OwnerID   *int32  `json:"ownerid"`
}

type UpdateCarResponse struct {
	ID int32 `json:"id"`
}

type CreateCarrierRequest struct {
	ID             int32  `json:"id"`
	DriverCategory string `json:"drivercategory"`
}

type CreateCarrierResponse struct {
	ID int32 `json:"id"`
}

type GetCarrierResponse struct {
	ID             int32  `json:"id"`
	DriverCategory string `json:"drivercategory"`
}

type UpdateCarrierRequest struct {
	DriverCategory string `json:"drivercategory"`
}

type UpdateCarrierResponse struct {
	ID int32 `json:"id"`
}

type CreateRecipientRequest struct {
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	ThirdName  string `json:"thirdname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type CreateRecipientResponse struct {
	ID int32 `json:"id"`
}

type GetRecipientResponse struct {
	ID         int32  `json:"id"`
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	ThirdName  string `json:"thirdname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type UpdateRecipientRequest struct {
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	ThirdName  string `json:"thirdname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type UpdateRecipientResponse struct {
	ID int32 `json:"id"`
}