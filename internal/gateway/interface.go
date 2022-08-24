package gateway

import (
	"context"
	"errors"
)

type DAO interface {
	GetByIP(context.Context, string) (*GeoLocation, error)
}

type GeoLocation struct {
	IPAddress    string  `json:"ip_address"`
	CountryCode  string  `json:"country_code"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	MysteryValue int64   `json:"mystery_value"`
}

var ErrNoData = errors.New("no data")
