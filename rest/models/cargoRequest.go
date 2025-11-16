package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	StatusWaitingTripChoice     = "WAITING_TRIP_CHOICE"
	StatusWaitingDriverApproval = "WAITING_DRIVER_APPROVAL"
	StatusApprovedByDriver      = "APPROVED_BY_DRIVER"
	StatusRejectedByDriver      = "REJECTED_BY_DRIVER"
	StatusInDelivery            = "IN_DELIVERY"
	StatusCompleted             = "COMPLETED"
	StatusCanceled              = "CANCELED"
)

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
