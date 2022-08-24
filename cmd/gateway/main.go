package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/nimusp/geolocation_service/internal/gateway"
	"github.com/nimusp/geolocation_service/internal/storage"
	"github.com/nimusp/geolocation_service/internal/storage/database"
)

const (
	portEnvName    = "GATEWAY_PORT"
	timeoutEnvName = "GATEWAY_TIMEOUT_SECONDS"

	defaultPort    = "8888"
	defaultTimeout = 5

	locationHandlerPath = "/location/{ipAddress}"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	db, err := database.New()
	if err != nil {
		sugaredLogger.Fatalf("init DB conn error: %s", err)
	}

	storage := storage.NewStorage(db)
	api := gateway.New(storage, sugaredLogger)

	r := mux.NewRouter()
	r.HandleFunc(locationHandlerPath, api.GeoLocationHandler).Methods(http.MethodGet)

	port := parseEnvPort()
	timeout := parseEnvTimeout()

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: timeout * time.Second,
		ReadTimeout:  timeout * time.Second,
	}

	sugaredLogger.Infof("start listening :%s", port)

	log.Fatal(srv.ListenAndServe())
}

func parseEnvPort() string {
	port := os.Getenv(portEnvName)
	if len(port) == 0 {
		port = defaultPort
	}
	return port
}

func parseEnvTimeout() time.Duration {
	timeoutEnv := os.Getenv(timeoutEnvName)
	if len(timeoutEnv) == 0 {
		return defaultTimeout
	}

	timeout, err := strconv.ParseUint(timeoutEnv, 10, 64)
	if err != nil {
		return defaultTimeout
	}

	if timeout <= 0 {
		return defaultTimeout
	}

	return time.Duration(timeout)
}
