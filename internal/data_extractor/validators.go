package data_exctractor

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func trimRow(row []string) (ip, countryCode, country, city, lat, lon, value string) {
	ip = strings.TrimSpace(row[ipAddressPosition])
	countryCode = strings.TrimSpace(row[countryCodePosition])
	country = strings.TrimSpace(row[countryPosition])
	city = strings.TrimSpace(row[cityPosition])
	lat = strings.TrimSpace(row[latitudePosition])
	lon = strings.TrimSpace(row[longitudePosition])
	value = strings.TrimSpace(row[mysteryValuePosition])
	return
}

var (
	countryCodeRE    = regexp.MustCompile(`^[A-Z]{2}$`) //  Alpha-2 code
	countryAndCityRE = regexp.MustCompile(`^[A-Z]{1}[a-z ]+`)
)

const (
	minLatitude  = -90.
	maxLatitude  = 90.
	minLongitude = -180.
	maxLongitude = 180.
)

func isValidIP(ipAddress string) bool {
	if len(ipAddress) == 0 {
		return false
	}

	parsedAddr := net.ParseIP(ipAddress)
	return parsedAddr != nil
}

func isValidCountryCode(countryCode string) bool {
	return countryCodeRE.MatchString(countryCode)
}

func isValidCountry(country string) bool {
	return countryAndCityRE.MatchString(country)
}

func isValidCity(city string) bool {
	return countryAndCityRE.MatchString(city)
}

func parseLatitude(rawLat string) (float64, error) {
	val, err := strconv.ParseFloat(rawLat, 64)
	if err != nil {
		return 0, err
	}

	if val < minLatitude || val > maxLatitude {
		return 0, fmt.Errorf("invalid latitude value %v", val)
	}

	return val, nil
}

func parseLongitude(rawLon string) (float64, error) {
	val, err := strconv.ParseFloat(rawLon, 64)
	if err != nil {
		return 0, err
	}

	if val < minLongitude || val > maxLongitude {
		return 0, fmt.Errorf("invalid longitude value %v", val)
	}

	return val, nil
}

func parseMysteryValue(rawValue string) (int64, error) {
	return strconv.ParseInt(rawValue, 10, 64)
}
