package models

type (
	Coords struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	Station struct {
		Address string `json:"address"`
		Coords  Coords `json:"coords"`
	}
)
