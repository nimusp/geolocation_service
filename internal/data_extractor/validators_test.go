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

func Test_isValidIP(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "empty IP",
			ip:   "",
			want: false,
		},
		{
			name: "bad IP",
			ip:   "bad IP",
			want: false,
		},
		{
			name: "IP v4",
			ip:   "127.0.0.1",
			want: true,
		},
		{
			name: "IP v6",
			ip:   "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidIP(tc.ip)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_isValidCountryCode(t *testing.T) {
	testCases := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "empty code",
			code: "",
			want: false,
		},
		{
			name: "invalid code 1",
			code: "A",
			want: false,
		},
		{
			name: "invalid code 2",
			code: "aa",
			want: false,
		},
		{
			name: "Alpha 3 format",
			code: "NLD",
			want: false,
		},
		{
			name: "falid",
			code: "NL",
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidCountryCode(tc.code)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_isValidCountry(t *testing.T) {
	testCases := []struct {
		name    string
		country string
		want    bool
	}{
		{
			name:    "empty",
			country: "",
			want:    false,
		},
		{
			name:    "first latter is lowercased",
			country: "netherlands",
			want:    false,
		},
		{
			name:    "with space",
			country: "United Kingdom of Great Britain and Northern Ireland",
			want:    true,
		},
		{
			name:    "valid",
			country: "Netherlands",
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidCountry(tc.country)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_isValidCity(t *testing.T) {
	testCases := []struct {
		name string
		city string
		want bool
	}{
		{
			name: "empty",
			city: "",
			want: false,
		},
		{
			name: "first latter is lowercased",
			city: "amsterdam",
			want: false,
		},
		{
			name: "with space",
			city: "Kuala Lumpur",
			want: true,
		},
		{
			name: "valid",
			city: "Amsterdam",
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidCity(tc.city)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_parseLatitude(t *testing.T) {
	testCases := []struct {
		name     string
		latitude string
		want     float64
		wantErr  bool
	}{
		{
			name:     "empty",
			latitude: "",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "not float",
			latitude: "hello",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "invalid value 1",
			latitude: "-91",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "invalid value 2",
			latitude: "91",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "valid",
			latitude: "-1.234",
			want:     -1.234,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := parseLatitude(tc.latitude)

			assert.Equal(t, tc.wantErr, gotErr != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_parseLongitude(t *testing.T) {
	testCases := []struct {
		name      string
		longitude string
		want      float64
		wantErr   bool
	}{
		{
			name:      "empty",
			longitude: "",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "not float",
			longitude: "hello",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "invalid value 1",
			longitude: "-182",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "invalid value 2",
			longitude: "182",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "valid",
			longitude: "11",
			want:      11,
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := parseLongitude(tc.longitude)

			assert.Equal(t, tc.wantErr, gotErr != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_parseMysteryValue(t *testing.T) {
	testCases := []struct {
		name    string
		value   string
		want    int64
		wantErr bool
	}{
		{
			name:    "empty",
			value:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "not int",
			value:   "hello",
			want:    0,
			wantErr: true,
		},
		{
			name:    "valid",
			value:   "111",
			want:    111,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := parseMysteryValue(tc.value)

			assert.Equal(t, tc.wantErr, gotErr != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
