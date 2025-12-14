package models

import "github.com/google/uuid"

type (
	GetTripRouteResponse struct {
		ID       uuid.UUID          `json:"id"`
		Stations []TripStationPoint `json:"stations"`
		Points   []float64          `json:"points"`
	}

	TripStationPoint struct {
		Station       Station `json:"station"`
		Distance      int     `json:"distance"`
		OrderNum      int     `json:"orderNum"`
		ArrivalAt     int64   `json:"arrivalAt"`
		DepartureTime int64   `json:"departureTime"`
	}
)
