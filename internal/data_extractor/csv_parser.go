package data_exctractor

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
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

type csvParser struct{}

func (csvParser) Extract(path string) (*Data, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s: %w", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read file %s error: %w", path, err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("too few rows (%d)", len(rows))
	}

	data := &Data{Rows: make([]Row, 0, len(rows))}
	for i := 1; i < len(rows); i++ { // skip header
		if len(rows[i]) != columsNumber {
			continue
		}

		var isBadIP, isBadCountryCode, isBadCountry, isBadCity, isBadLat, isBadLon, isBadValue bool

		if len(rows[i][ipAddressPosition]) == 0 {
			isBadIP = true
			data.Stats.BadIPAddress++
		}

		if len(rows[i][countryCodePosition]) == 0 {
			isBadCountryCode = true
			data.Stats.BadCountryCode++
		}

		if len(rows[i][countryPosition]) == 0 {
			isBadCountry = true
			data.Stats.BadCountry++
		}

		if len(rows[i][cityPosition]) == 0 {
			isBadCity = true
			data.Stats.BadCity++
		}

		lat, err := strconv.ParseFloat(rows[i][latitudePosition], 64)
		if err != nil {
			isBadLat = true
			data.Stats.BadLatitude++
		}

		lon, err := strconv.ParseFloat(rows[i][longitudePosition], 64)
		if err != nil {
			isBadLon = true
			data.Stats.BadLongitude++
		}

		value, err := strconv.ParseInt(rows[i][mysteryValuePosition], 10, 64)
		if err != nil {
			isBadValue = true
			data.Stats.BadMysteryValue++
		}

		if isBadIP || isBadCountryCode || isBadCountry || isBadCity || isBadLat || isBadLon || isBadValue {
			continue
		}

		data.Rows = append(data.Rows, Row{
			IPAddress:    rows[i][ipAddressPosition],
			CountryCode:  rows[i][countryCodePosition],
			Country:      rows[i][countryPosition],
			City:         rows[i][cityPosition],
			Latitude:     lat,
			Longitude:    lon,
			MysteryValue: value,
		})
	}

	if len(data.Rows) == 0 {
		return nil, fmt.Errorf("no valid rows in %d raw rows", len(rows))
	}

	return data, nil
}
