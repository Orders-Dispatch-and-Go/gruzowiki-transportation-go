package models

import "github.com/google/uuid"

const (
	TripStatusPending    = "PENDING"
	TripStatusInProgress = "IN_PROGRESS"
	TripStatusCompleted  = "COMPLETED"
	TripStatusCanceled   = "CANCELED"
)

type (
	TripResponse struct {
		ID              uuid.UUID  `json:"id"`
		RouteID         *uuid.UUID `json:"route_id"`
		FromStation     Station    `json:"fromStation"`
		ToStation       Station    `json:"toStation"`
		StartedAt       int64      `json:"startedAt"`
		CalculatedEndAt int64      `json:"calculatedEndAt"`
		ActualEndAt     int64      `json:"actualEndAt"`
		Price           string     `json:"price"`
		Status          string     `json:"status"`
		CarrierID       int        `json:"carrierId"`
		CarID           int        `json:"carId"`
	}

	Trips struct {
		Trips []TripResponse `json:"trips"`
	}

	GetTripResponse struct {
		ID             uuid.UUID
		RouteID        *uuid.UUID
		FromStation    *uuid.UUID
		ToStation      *uuid.UUID
		StartedAt      int64
		CalculateEndAt int64
		ActualEndAt    int64
		Price          string
		Status         string
		Carrier        *int32
		Car            *int32
	}
)

type CreateTripRequest struct {
	Carrier int32 `json:"carrier"`

	FromStation Station `json:"fromStation"`
	ToStation   Station `json:"toStation"`

	StartedAt string `json:"startedAt"`
}

type CreateTripResponse struct {
	ID uuid.UUID `json:"id"`
}
