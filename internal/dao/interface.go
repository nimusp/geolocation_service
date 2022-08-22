package dao

import (
	"context"
	"errors"
)

type DAO interface {
	GetByIP(context.Context, string) (*GeoLocation, error)
}

type GeoLocation struct {
	IPAddress    string
	CountryCode  string
	Country      string
	City         string
	Latitude     float64
	Longitude    float64
	MysteryValue int64
}

var ErrNoData = errors.New("no data")
