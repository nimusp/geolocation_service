package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/nimusp/geolocation_service/internal/storage"
)

func Test_Insert(t *testing.T) {
	someErr := errors.New("some")
	val := storage.GeoLocation{
		IPAddress:    "0.0.0.0",
		CountryCode:  "NL",
		Country:      "Netherlands",
		City:         "Amsterdam",
		Latitude:     52.377956,
		Longitude:    4.897070,
		MysteryValue: 42,
	}

	testCases := []struct {
		name      string
		beginErr  error
		insertErr error
		commitErr error
		wantErr   bool
	}{
		{
			name:     "begin err",
			beginErr: someErr,
			wantErr:  true,
		},
		{
			name:      "insert err",
			insertErr: someErr,
			wantErr:   true,
		},
		{
			name:      "success",
			insertErr: nil,
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockConn, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockConn.Close()

			sqlMock.ExpectBegin().WillReturnError(tc.beginErr)

			if tc.beginErr == nil {
				sqlMock.
					ExpectExec("^INSERT INTO geolocation (.+)$").
					WithArgs(val.IPAddress, val.CountryCode, val.Country, val.City, val.Latitude, val.Longitude, val.MysteryValue).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(tc.insertErr)

				if tc.insertErr != nil {
					sqlMock.ExpectCommit()
				} else {
					sqlMock.ExpectRollback()
				}

			}

			d := db{conn: mockConn}

			gotErr := d.Insert(context.Background(), []storage.GeoLocation{val}, 1)

			assert.Equal(t, tc.wantErr, gotErr != nil)
		})
	}
}

func Test_Select(t *testing.T) {
	testIP := "0.0.0.0"
	commonErr := errors.New("some")

	testCases := []struct {
		name    string
		rows    []driver.Value
		rowsErr error
		dbErr   error
		want    *storage.GeoLocation
		wantErr error
	}{
		{
			name:    "no data",
			rows:    nil,
			rowsErr: sql.ErrNoRows,
			want:    nil,
			wantErr: storage.ErrNoData,
		},
		{
			name:    "db error",
			rows:    nil,
			dbErr:   commonErr,
			want:    nil,
			wantErr: commonErr,
		},
		{
			name:    "no data",
			rows:    []driver.Value{testIP, "NL", "Netherlands", "Amsterdam", 52.377956, 4.897070, 42},
			rowsErr: nil,
			want: &storage.GeoLocation{
				IPAddress:    testIP,
				CountryCode:  "NL",
				Country:      "Netherlands",
				City:         "Amsterdam",
				Latitude:     52.377956,
				Longitude:    4.897070,
				MysteryValue: 42,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockConn, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockConn.Close()

			columns := []string{
				"ip_address",
				"country_code",
				"country",
				"city",
				"latitude",
				"longitude",
				"mystery_value",
			}

			mockedRows := sqlmock.NewRows(columns).RowError(0, tc.rowsErr)
			if len(tc.rows) > 0 {
				mockedRows = mockedRows.AddRow(tc.rows...)
			}

			sqlMock.ExpectQuery("^SELECT (.+) FROM geolocation").
				WithArgs(testIP).
				WillReturnRows(mockedRows).
				WillReturnError(tc.dbErr)

			d := db{conn: mockConn}

			got, gotErr := d.Select(context.Background(), testIP)

			if tc.wantErr != nil {
				if errors.Is(tc.wantErr, storage.ErrNoData) {
					assert.True(t, errors.Is(gotErr, tc.wantErr))
				} else {
					assert.True(t, errors.Is(gotErr, commonErr))
				}
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_splitToChunks(t *testing.T) {
	testCases := []struct {
		name string
		data []storage.GeoLocation
		size uint
		want [][]storage.GeoLocation
	}{
		{
			name: "single chunk",
			data: []storage.GeoLocation{{IPAddress: "0.0.0.0"}, {IPAddress: "1.1.1.1"}},
			size: 3,
			want: [][]storage.GeoLocation{{{IPAddress: "0.0.0.0"}, {IPAddress: "1.1.1.1"}}},
		},
		{
			name: "odd",
			data: []storage.GeoLocation{{IPAddress: "0.0.0.0"}, {IPAddress: "1.1.1.1"}, {IPAddress: "2.2.2.2"}},
			size: 2,
			want: [][]storage.GeoLocation{{{IPAddress: "0.0.0.0"}, {IPAddress: "1.1.1.1"}}, {{IPAddress: "2.2.2.2"}}},
		},
		{
			name: "even",
			data: []storage.GeoLocation{{IPAddress: "0.0.0.0"}, {IPAddress: "1.1.1.1"}, {IPAddress: "2.2.2.2"}},
			size: 1,
			want: [][]storage.GeoLocation{{{IPAddress: "0.0.0.0"}}, {{IPAddress: "1.1.1.1"}}, {{IPAddress: "2.2.2.2"}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := splitToChunks(tc.data, tc.size)

			assert.Equal(t, tc.want, got)
		})
	}
}
