package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	CargoRequestStatusPending    = "PENDING"
	CargoRequestStatusInProgress = "IN_PROGRESS"
	CargoRequestStatus           = "COMPLETED"
	CargoRequestStatusCanceled   = "CANCELED"
)

type (
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
		ID          uuid.UUID `json:"id"`
		ReceiveCode string    `json:"receiveCode"`
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
		ReceiveCode *string  `json:"receiveCode"`
		Worth       int      `json:"worth"`
		Width       int      `json:"width"`
		Height      int      `json:"height"`
		Length      int      `json:"length"`
		Weight      int      `json:"weight"`
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

	CargoRequestIdsRequest struct {
		Ids []string `json:"cargoRequests"`
	}

	PotentialRoutesRequest struct {
		TripRouteId          string   `json:"tripRouteId"`
		CargoRequestRouteIds []string `json:"cargoRequestRouteIds"`
	}

	PotentialRoutesResponse struct {
		CargoRequests []string `json:"cargoRequests"`
	}

	GetCargoRequestsForTripResponse struct {
		CargoRequests []CargoRequestResponse `json:"cargoRequests"`
	}

	GetCargoRequestsForTripFilter struct {
		CargoLengthMax *int32
		CargoWidthMax  *int32
		CargoHeightMax *int32
		CargoType      *int32
		Deadline       *int64
		MinPrice       *int64
	}
)
