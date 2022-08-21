package data_exctractor

import (
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
