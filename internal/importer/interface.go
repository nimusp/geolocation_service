package importer

import "context"

type DataSaver interface {
	Save(context.Context, []GeoLocation, uint) error
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

type Exctactor interface {
	Extract() (*Data, error)
}

type Data struct {
	Rows  []GeoLocation
	Stats Statics
}

type Statics struct {
	RawRows         int64
	Doubles         int64
	BadIPAddress    int64
	BadCountryCode  int64
	BadCountry      int64
	BadCity         int64
	BadLatitude     int64
	BadLongitude    int64
	BadMysteryValue int64
}
