package storage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nimusp/geolocation_service/internal/gateway"
	"github.com/nimusp/geolocation_service/internal/importer"
	"github.com/nimusp/geolocation_service/internal/storage"
	mock_database "github.com/nimusp/geolocation_service/mocks/database"
)

func Test_Save(t *testing.T) {
	testCases := []struct {
		name      string
		insertErr error
		wantErr   bool
	}{
		{
			name:      "insert err",
			insertErr: errors.New("some"),
			wantErr:   true,
		},
		{
			name:      "no error",
			insertErr: nil,
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := mock_database.NewDatabase(t)
			mockDB.EXPECT().Insert(mock.Anything, mock.Anything, mock.Anything).Return(tc.insertErr)

			s := storage.NewStorage(mockDB)

			gotErr := s.Save(context.Background(), []importer.GeoLocation{{IPAddress: "0.0.0.0"}}, 1)

			assert.Equal(t, tc.wantErr, gotErr != nil)
		})
	}
}

func Test_GetByIP(t *testing.T) {
	testCases := []struct {
		name      string
		selectRes *storage.GeoLocation
		selectErr error
		want      *gateway.GeoLocation
		wantErr   bool
	}{
		{
			name:      "select err",
			selectRes: nil,
			selectErr: errors.New("some"),
			wantErr:   true,
		},
		{
			name:      "no data",
			selectRes: nil,
			selectErr: storage.ErrNoData,
			wantErr:   true,
		},
		{
			name: "no error",
			selectRes: &storage.GeoLocation{
				IPAddress:    "0.0.0.0",
				CountryCode:  "NL",
				Country:      "Netherlands",
				City:         "Amsterdam",
				Latitude:     52.377956,
				Longitude:    4.897070,
				MysteryValue: 42,
			},
			selectErr: nil,
			want: &gateway.GeoLocation{
				IPAddress:    "0.0.0.0",
				CountryCode:  "NL",
				Country:      "Netherlands",
				City:         "Amsterdam",
				Latitude:     52.377956,
				Longitude:    4.897070,
				MysteryValue: 42,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := mock_database.NewDatabase(t)
			mockDB.EXPECT().Select(mock.Anything, mock.Anything).Return(tc.selectRes, tc.selectErr)

			s := storage.NewStorage(mockDB)

			got, gotErr := s.GetByIP(context.Background(), "0.0.0.0")

			assert.Equal(t, tc.wantErr, gotErr != nil)
			if tc.selectErr != nil && errors.Is(tc.selectErr, storage.ErrNoData) {
				assert.True(t, errors.Is(gotErr, gateway.ErrNoData))
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
