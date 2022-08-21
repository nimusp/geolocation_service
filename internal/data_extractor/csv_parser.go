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
		ip, countryCode, country, city, latRaw, lonRaw, valueRaw := trimRow(rows[i])

		if len(ip) == 0 {
			isBadIP = true
			data.Stats.BadIPAddress++
		}

		if len(countryCode) == 0 {
			isBadCountryCode = true
			data.Stats.BadCountryCode++
		}

		if len(country) == 0 {
			isBadCountry = true
			data.Stats.BadCountry++
		}

		if len(city) == 0 {
			isBadCity = true
			data.Stats.BadCity++
		}

		lat, err := strconv.ParseFloat(latRaw, 64)
		if err != nil {
			isBadLat = true
			data.Stats.BadLatitude++
		}

		lon, err := strconv.ParseFloat(lonRaw, 64)
		if err != nil {
			isBadLon = true
			data.Stats.BadLongitude++
		}

		value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			isBadValue = true
			data.Stats.BadMysteryValue++
		}

		if isBadIP || isBadCountryCode || isBadCountry || isBadCity || isBadLat || isBadLon || isBadValue {
			continue
		}

		data.Rows = append(data.Rows, Row{
			IPAddress:    ip,
			CountryCode:  countryCode,
			Country:      country,
			City:         city,
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
