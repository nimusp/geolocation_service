package storage

import (
	"context"
	"errors"
)

type Database interface {
	Insert(ctx context.Context, data []GeoLocation, batchSize uint) error
	Select(ctx context.Context, ipAddress string) (*GeoLocation, error)
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
