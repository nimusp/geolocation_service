package data_exctractor

type Exctactor interface {
	Extract(path string) (*Data, error)
}

type Data struct {
	Rows  []Row
	Stats Statics
}

type Statics struct {
	BadIPAddress    int64
	BadCountryCode  int64
	BadCountry      int64
	BadCity         int64
	BadLatitude     int64
	BadLongitude    int64
	BadMysteryValue int64
}

type Row struct {
	IPAddress    string
	CountryCode  string
	Country      string
	City         string
	Latitude     float64
	Longitude    float64
	MysteryValue int64
}
