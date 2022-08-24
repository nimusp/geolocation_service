package gateway_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/nimusp/geolocation_service/internal/gateway"
	mock_gateway "github.com/nimusp/geolocation_service/mocks/gateway"
)

func Test_GeoLocationHandler(t *testing.T) {
	testCases := []struct {
		name           string
		ipAddress      string
		daoResp        *gateway.GeoLocation
		daoErr         error
		wantStatusCode int
		wantResp       gateway.GeoLocation
	}{
		{
			name:           "bad ip",
			ipAddress:      "127.0",
			daoResp:        nil,
			daoErr:         nil,
			wantStatusCode: http.StatusBadRequest,
			wantResp:       gateway.GeoLocation{},
		},
		{
			name:           "ip not found",
			ipAddress:      "127.0.0.1",
			daoResp:        nil,
			daoErr:         gateway.ErrNoData,
			wantStatusCode: http.StatusNotFound,
			wantResp:       gateway.GeoLocation{},
		},
		{
			name:           "internal error",
			ipAddress:      "127.0.0.1",
			daoResp:        nil,
			daoErr:         errors.New("some"),
			wantStatusCode: http.StatusInternalServerError,
			wantResp:       gateway.GeoLocation{},
		},
		{
			name:      "http 200",
			ipAddress: "0.0.0.0",
			daoResp: &gateway.GeoLocation{
				IPAddress:    "0.0.0.0",
				CountryCode:  "NL",
				Country:      "Netherlands",
				City:         "Amsterdam",
				Latitude:     52.377956,
				Longitude:    4.897070,
				MysteryValue: 42,
			},
			daoErr:         nil,
			wantStatusCode: http.StatusOK,
			wantResp: gateway.GeoLocation{
				IPAddress:    "0.0.0.0",
				CountryCode:  "NL",
				Country:      "Netherlands",
				City:         "Amsterdam",
				Latitude:     52.377956,
				Longitude:    4.897070,
				MysteryValue: 42,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDAO := new(mock_gateway.DAO)
			mockDAO.Test(t)
			mockDAO.EXPECT().GetByIP(mock.Anything, mock.Anything).Return(tc.daoResp, tc.daoErr)

			a := gateway.New(mockDAO, zap.NewNop().Sugar())

			path := fmt.Sprintf("/location/%s", tc.ipAddress)
			req, err := http.NewRequest(http.MethodGet, path, nil)
			if err != nil {
				t.Fatal(err)
			}

			router := mux.NewRouter()
			router.HandleFunc("/location/{ipAddress}", a.GeoLocationHandler)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			res := w.Result()

			rawResp, err := io.ReadAll(res.Body)
			if err != nil {
				res.Body.Close()
				t.Errorf("expected error to be nil got %v", err)
			}
			res.Body.Close()

			assert.Equal(t, tc.wantStatusCode, res.StatusCode)
			if tc.wantStatusCode != http.StatusOK {
				return
			}

			var response gateway.GeoLocation
			err = json.Unmarshal(rawResp, &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, response)
		})
	}

}
