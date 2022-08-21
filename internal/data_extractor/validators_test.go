package data_exctractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trimRow(t *testing.T) {
	ip := "200.106.141.15"
	countryCode := "SI"
	country := "Nepal"
	city := "DuBuquemouth"
	lat := "-84.87503094689836"
	lon := "7.206435933364332"
	value := "7823011346"

	testCases := []struct {
		name string
		row  []string
	}{
		{
			name: "ip with space",
			row:  []string{" 200.106.141.15 ", "SI", "Nepal", "DuBuquemouth", "-84.87503094689836", "7.206435933364332", "7823011346"},
		},
		{
			name: "country code with space",
			row:  []string{"200.106.141.15", " SI ", "Nepal", "DuBuquemouth", "-84.87503094689836", "7.206435933364332", "7823011346"},
		},
		{
			name: "country space",
			row:  []string{"200.106.141.15", "SI", " Nepal ", "DuBuquemouth", "-84.87503094689836", "7.206435933364332", "7823011346"},
		},
		{
			name: "city with space",
			row:  []string{"200.106.141.15", "SI", "Nepal", " DuBuquemouth ", "-84.87503094689836", "7.206435933364332", "7823011346"},
		},
		{
			name: "latitude with space",
			row:  []string{"200.106.141.15", "SI", "Nepal", "DuBuquemouth", " -84.87503094689836 ", "7.206435933364332", "7823011346"},
		},
		{
			name: "longitude with space",
			row:  []string{"200.106.141.15", "SI", "Nepal", "DuBuquemouth", "-84.87503094689836", " 7.206435933364332 ", "7823011346"},
		},
		{
			name: "mustery valie with space",
			row:  []string{"200.106.141.15", "SI", "Nepal", "DuBuquemouth", "-84.87503094689836", "7.206435933364332", " 7823011346 "},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotIP, gotCountryCode, gotCountry, gotCity, gotLat, gotLon, gotvVlue := trimRow(tc.row)

			assert.Equal(t, ip, gotIP)
			assert.Equal(t, countryCode, gotCountryCode)
			assert.Equal(t, country, gotCountry)
			assert.Equal(t, city, gotCity)
			assert.Equal(t, lat, gotLat)
			assert.Equal(t, lon, gotLon)
			assert.Equal(t, value, gotvVlue)
		})
	}
}
