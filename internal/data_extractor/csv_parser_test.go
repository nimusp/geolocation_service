package data_exctractor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nimusp/geolocation_service/internal/importer"
)

func Test_Extract(t *testing.T) {
	testCases := []struct {
		name    string
		data    []byte
		want    *importer.Data
		wantErr bool
	}{
		{
			name: "bad data format",
			data: []byte(
				`ip_address country_code country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115`,
			),
			want:    nil,
			wantErr: true,
		},
		{
			name: "too few rows",
			data: []byte(
				`ip_address,country_code,country,city,latitude,longitude,mystery_value`,
			),
			want:    nil,
			wantErr: true,
		},
		{
			name: "with bad rows",
			data: []byte(
				`ip_address,country_code,country,city,latitude,longitude,mystery_value
,,,,,abc,cde
,,,,-abc,,
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346`,
			),
			want: &importer.Data{
				Rows: []importer.Row{
					{
						IPAddress:    "200.106.141.15",
						CountryCode:  "SI",
						Country:      "Nepal",
						City:         "DuBuquemouth",
						Latitude:     -84.87503094689836,
						Longitude:    7.206435933364332,
						MysteryValue: 7823011346,
					},
				},
				Stats: importer.Statics{
					RawRows:         4,
					BadIPAddress:    2,
					BadCountryCode:  2,
					BadCountry:      2,
					BadCity:         2,
					BadLatitude:     2,
					BadLongitude:    2,
					BadMysteryValue: 2,
				},
			},
			wantErr: false,
		},
		{
			name: "no valid rows",
			data: []byte(
				`ip_address,country_code,country,city,latitude,longitude,mystery_value
,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,,New Neva,-68.31023296602508,-37.62435199624531,7301823115`,
			),
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid data",
			data: []byte(
				`ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115`,
			),
			want: &importer.Data{
				Rows: []importer.Row{
					{
						IPAddress:    "200.106.141.15",
						CountryCode:  "SI",
						Country:      "Nepal",
						City:         "DuBuquemouth",
						Latitude:     -84.87503094689836,
						Longitude:    7.206435933364332,
						MysteryValue: 7823011346,
					},
					{
						IPAddress:    "160.103.7.140",
						CountryCode:  "CZ",
						Country:      "Nicaragua",
						City:         "New Neva",
						Latitude:     -68.31023296602508,
						Longitude:    -37.62435199624531,
						MysteryValue: 7301823115,
					},
				},
				Stats: importer.Statics{RawRows: 3},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dirName := t.TempDir()
			path := filepath.Join(dirName, "data.csv")
			err := os.WriteFile(path, tc.data, os.ModePerm)
			assert.NoError(t, err)

			parser := csvParser{path: path}
			got, gotErr := parser.Extract()

			assert.Equal(t, tc.wantErr, gotErr != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_Extract_badFilePath(t *testing.T) {
	dirName := t.TempDir()
	path := filepath.Join(dirName, "data.csv")

	parser := csvParser{path: path}
	_, err := parser.Extract()

	assert.Error(t, err)
}

func Test_distinctByIP(t *testing.T) {
	testCases := []struct {
		name string
		rows []importer.Row
		want []importer.Row
	}{
		{
			name: "with doubles",
			rows: []importer.Row{{IPAddress: "127.0.0.1"}, {IPAddress: "8.8.8.8"}, {IPAddress: "127.0.0.1"}},
			want: []importer.Row{{IPAddress: "127.0.0.1"}, {IPAddress: "8.8.8.8"}},
		},
		{
			name: "without doubles",
			rows: []importer.Row{{IPAddress: "127.0.0.1"}, {IPAddress: "8.8.8.8"}, {IPAddress: "192.168.0.1"}},
			want: []importer.Row{{IPAddress: "127.0.0.1"}, {IPAddress: "8.8.8.8"}, {IPAddress: "192.168.0.1"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := distinctByIP(tc.rows)

			assert.Equal(t, tc.want, got)
		})
	}
}
