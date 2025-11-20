package models

type CreateRecipientRequest struct {
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

type UpdateRecipientResponse struct {
	ID int32 `json:"id"`
}