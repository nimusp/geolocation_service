package main

import (
	"context"
	"os"
	"strconv"

	"go.uber.org/zap"

	data_exctractor "github.com/nimusp/geolocation_service/internal/data_extractor"
	"github.com/nimusp/geolocation_service/internal/importer"
	"github.com/nimusp/geolocation_service/internal/storage"
	"github.com/nimusp/geolocation_service/internal/storage/database"
)

const (
	importPathEnvName      = "IMPORT_FROM"
	importBatchSizeEnvName = "IMPORT_BATCH_SIZE"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	importPath := os.Getenv(importPathEnvName)
	if len(importPath) == 0 {
		sugaredLogger.Fatalln("empty import path")
	}

	db, err := database.New()
	if err != nil {
		sugaredLogger.Fatalf("init DB conn error: %s", err)
	}

	storage := storage.NewStorage(db)
	extractor := data_exctractor.New(importPath)

	customBatchSize, _ := strconv.ParseUint(os.Getenv(importBatchSizeEnvName), 10, 64)
	importer := importer.New(storage, extractor, sugaredLogger, importer.WithBatchSize(uint(customBatchSize)))
	if err := importer.Import(context.Background()); err != nil {
		sugaredLogger.Fatalln("import error", err)
	}
}
