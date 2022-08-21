package data_exctractor

import (
	"encoding/csv"
	"fmt"
	"os"
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

func New() *csvParser { return new(csvParser) }

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

	data := &Data{
		Rows:  make([]Row, 0, len(rows)),
		Stats: Statics{RawRows: int64(len(rows))},
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

	unique := distinctByIP(data.Rows)
	data.Stats.Doubles = int64(len(data.Rows) - len(unique))
	data.Rows = unique

	if len(data.Rows) == 0 {
		return nil, fmt.Errorf("no valid rows in %d raw rows", len(rows))
	}

	return data, nil
}

func distinctByIP(rows []Row) []Row {
	set := make(map[string]struct{}, len(rows))
	unique := make([]Row, 0, len(rows))

	for _, row := range rows {
		if _, exists := set[row.IPAddress]; !exists {
			set[row.IPAddress] = struct{}{}
			unique = append(unique, row)
		}
	}

	return unique
}
