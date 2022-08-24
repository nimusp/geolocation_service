package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type api struct {
	dao    DAO
	logger *zap.SugaredLogger
}

func New(dao DAO, logger *zap.SugaredLogger) *api {
	return &api{dao: dao, logger: logger}
}

func (a *api) GeoLocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ipAddress := vars["ipAddress"]

	if net.ParseIP(ipAddress) == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad IP address: %s", ipAddress)
		return
	}

	info, err := a.dao.GetByIP(r.Context(), ipAddress)
	if err != nil {
		if errors.Is(err, ErrNoData) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Unknown IP address: %s", ipAddress)
			a.logger.Warnf("got unknonw IP address %s", ipAddress)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		a.logger.Errorf("handle IP %s error: %s", ipAddress, err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(info); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.logger.Errorf("marshall data with IP %s error: %s", ipAddress, err.Error())
		return
	}
}
