package client

import (
	"github.com/google/uuid"
	"gruzowiki/rest/models"
)

type (
	StationDTO struct {
		ID      uuid.UUID     `json:"id"`
		Address string        `json:"address"`
		Coords  StationCoords `json:"coords"`
	}

	StationCoords struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	CreateRouteRequestBody struct {
		FromStation StationDTO `json:"fromStation"`
		ToStation   StationDTO `json:"toStation"`
	}

	CreateRouteResponseBody struct {
		ID uuid.UUID `json:"id"`
	}

	MergeCargoRequestIntoTripRoute struct {
		TripRouteID         string   `json:"tripRouteId"`
		CargoRequestRouteID []string `json:"cargoRequestRouteId"`
	}

	GetTripRouteResponse struct {
		ID       uuid.UUID          `json:"id"`
		Stations []TripStationPoint `json:"stations"`
	}

	TripStationPoint struct {
		Station       models.Station `json:"station"`
		Distance      int            `json:"distance"`
		OrderNum      int            `json:"orderNum"`
		ArrivalAt     int64          `json:"arrivalAt"`
		DepartureTime int64          `json:"departureTime"`
	}

	GetRoutePointsResponse struct {
		Points [][]float64 `json:"points"`
	}
)
