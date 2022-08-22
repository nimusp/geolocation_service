package importer

import "context"

type Importer interface {
	Save(context.Context, []GeoLocation) error
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
