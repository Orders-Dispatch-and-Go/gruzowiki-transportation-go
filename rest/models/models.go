package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GetCarrierResponse struct {
	Id             int32  `json:"id"`
	DriverCategory string `json:"driverCategory"`
}

type Coords struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Station struct {
	Address string `json:"address"`
	Coords  Coords `json:"coords"`
}

type PostCargoRequestRequest struct {
	ConsignerID int32           `json:"consignerId"`
	RecipientID int32           `json:"recipientId"`
	FromStation Station         `json:"fromStation"`
	ToStation   Station         `json:"toStation"`
	Deadline    string          `json:"deadline"`
	MaxPrice    decimal.Decimal `json:"maxPrice"`
}

type PostCargoRequestResponse struct {
	ID uuid.UUID `json:"id"`
}
