package importer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/nimusp/geolocation_service/internal/importer"
	mock_importer "github.com/nimusp/geolocation_service/mocks/importer"
)

func Test_Import(t *testing.T) {
	testCases := []struct {
		name        string
		extractData *importer.Data
		extractErr  error
		saveErr     error
		saverCalls  int
		wantErr     bool
	}{
		{
			name:        "extract err",
			extractData: nil,
			extractErr:  errors.New("some"),
			saveErr:     nil,
			saverCalls:  0,
			wantErr:     true,
		},
		{
			name:        "save err",
			extractData: &importer.Data{Rows: make([]importer.GeoLocation, 0)},
			extractErr:  nil,
			saveErr:     errors.New("some"),
			saverCalls:  1,
			wantErr:     true,
		},
		{
			name:        "no err",
			extractData: &importer.Data{Rows: make([]importer.GeoLocation, 0)},
			extractErr:  nil,
			saveErr:     nil,
			saverCalls:  1,
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockExtractor := new(mock_importer.Exctactor)
			mockExtractor.Test(t)
			mockExtractor.EXPECT().Extract().Return(tc.extractData, tc.extractErr)

			mockSaver := new(mock_importer.DataSaver)
			mockSaver.Test(t)
			mockSaver.EXPECT().Save(mock.Anything, mock.Anything, mock.Anything).Return(tc.saveErr)

			i := importer.New(mockSaver, mockExtractor, zap.NewNop().Sugar())

			gotErr := i.Import(context.Background())

			assert.Equal(t, tc.wantErr, gotErr != nil)

			mockExtractor.AssertNumberOfCalls(t, "Extract", 1)
			mockSaver.AssertNumberOfCalls(t, "Save", tc.saverCalls)
		})
	}
}
