package data_exctractor

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/nimusp/geolocation_service/internal/importer"
)

// ip_address,country_code,country,city,latitude,longitude,mystery_value
const (
	ipAddressPosition = iota
	countryCodePosition
	countryPosition
	cityPosition
	latitudePosition
	longitudePosition
	mysteryValuePosition
)

const columsNumber = 7

type csvParser struct {
	path string
}

func New(path string) *csvParser {
	return &csvParser{path: path}
}

func (p *csvParser) Extract() (*importer.Data, error) {
	file, err := os.Open(p.path)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s: %w", p.path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read file %s error: %w", p.path, err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("too few rows (%d)", len(rows))
	}

	data := &importer.Data{
		Rows:  make([]importer.Row, 0, len(rows)),
		Stats: importer.Statics{RawRows: int64(len(rows))},
	}
	for i := 1; i < len(rows); i++ { // skip header
		if len(rows[i]) != columsNumber {
			continue
		}

		var isBadIP, isBadCountryCode, isBadCountry, isBadCity, isBadLat, isBadLon, isBadValue bool
		ip, countryCode, country, city, latRaw, lonRaw, valueRaw := trimRow(rows[i])

		if !isValidIP(ip) {
			isBadIP = true
			data.Stats.BadIPAddress++
		}

		if !isValidCountryCode(countryCode) {
			isBadCountryCode = true
			data.Stats.BadCountryCode++
		}

		if !isValidCountry(country) {
			isBadCountry = true
			data.Stats.BadCountry++
		}

		if !isValidCity(city) {
			isBadCity = true
			data.Stats.BadCity++
		}

		lat, err := parseLatitude(latRaw)
		if err != nil {
			isBadLat = true
			data.Stats.BadLatitude++
		}

		lon, err := parseLongitude(lonRaw)
		if err != nil {
			isBadLon = true
			data.Stats.BadLongitude++
		}

		value, err := parseMysteryValue(valueRaw)
		if err != nil {
			isBadValue = true
			data.Stats.BadMysteryValue++
		}

		if isBadIP || isBadCountryCode || isBadCountry || isBadCity || isBadLat || isBadLon || isBadValue {
			continue
		}

		data.Rows = append(data.Rows, importer.Row{
			IPAddress:    ip,
			CountryCode:  countryCode,
			Country:      country,
			City:         city,
			Latitude:     lat,
			Longitude:    lon,
			MysteryValue: value,
		})
	}

	unique := distinctByIP(data.Rows)
	data.Stats.Doubles = int64(len(data.Rows) - len(unique))
	data.Rows = unique

	if len(data.Rows) == 0 {
		return nil, fmt.Errorf("no valid rows in %d raw rows", len(rows))
	}

	return data, nil
}

func distinctByIP(rows []importer.Row) []importer.Row {
	set := make(map[string]struct{}, len(rows))
	unique := make([]importer.Row, 0, len(rows))

	for _, row := range rows {
		if _, exists := set[row.IPAddress]; !exists {
			set[row.IPAddress] = struct{}{}
			unique = append(unique, row)
		}
	}

	return unique
}
