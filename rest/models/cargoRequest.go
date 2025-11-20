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

type (
	Coords struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	Station struct {
		Address string `json:"address"`
		Coords  Coords `json:"coords"`
	}

	PostCargoRequestRequest struct {
		ConsignerID int32           `json:"consignerId"`
		RecipientID int32           `json:"recipientId"`
		FromStation Station         `json:"fromStation"`
		ToStation   Station         `json:"toStation"`
		Deadline    string          `json:"deadline"`
		MaxPrice    decimal.Decimal `json:"maxPrice"`
	}

	PostCargoRequestResponse struct {
		ID uuid.UUID `json:"id"`
	}

	CargoTypesResponse struct {
		CargoTypes []CargoType `json:"cargoTypes"`
	}

	CargoType struct {
		Id      int    `json:"id"`
		Type    string `json:"type"`
		Fragile bool   `json:"fragile"`
	}

	Cargo struct {
		Length         int    `json:"length"`
		Width          int    `json:"width"`
		Height         int    `json:"height"`
		Weight         int    `json:"weight"`
		CargoType      int    `json:"cargoType"`
		Description    string `json:"description"`
		Worth          int    `json:"worth"`
		CargoRequestId string `json:"cargoRequestId"`
	}

	CargoRequest struct {
		Cargo []Cargo `json:"cargo"`
	}

	IdsResponse struct {
		Ids []string `json:"ids"`
	}
)