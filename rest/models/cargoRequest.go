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

	GetStationsResponse struct {
		Stations []Station `json:"stations"`
	}

	CreateStationResponse struct {
		ID uuid.UUID `json:"id"`
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

	GetCargoRequest struct {
		ID          *string `json:"id"`
		ConsignerID *int    `json:"consignerId"`
		RecipientID *int    `json:"recipientId"`
		Status      *string `json:"status"`
		CreatedFrom *string `json:"createdFrom"`
		CreatedTo   *string `json:"createdTo"`
	}

	SearchCargoRequestsResponse struct {
		CargoRequests []CargoRequestResponse `json:"cargoRequests"`
	}

	CargoRequestResponse struct {
		ID          string   `json:"id"`
		ConsignerID int      `json:"consignerId"`
		RecipientID int      `json:"recipientId"`
		FromStation *Station `json:"fromStation"`
		ToStation   *Station `json:"toStation"`
		CreatedAt   int64    `json:"createdAt"`
		Deadline    int64    `json:"deadline"`
		RouteID     *string  `json:"routeId"`
		TripID      *string  `json:"tripId"`
		Price       string   `json:"price"`
		Status      string   `json:"status"`
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
