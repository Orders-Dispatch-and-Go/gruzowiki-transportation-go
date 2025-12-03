package client

import "github.com/google/uuid"

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
)
