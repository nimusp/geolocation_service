package importer

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

const defaultBatchSize = 1000

type dataImporter struct {
	saver     DataSaver
	extractor Exctactor
	batchSize uint
	logger    *zap.SugaredLogger
}

type opt func(*dataImporter)

func WithBatchSize(size uint) opt {
	return func(i *dataImporter) {
		if size > 0 {
			i.batchSize = size
		}
	}
}

func New(saver DataSaver, extractor Exctactor, logger *zap.SugaredLogger, opts ...opt) *dataImporter {
	dImporter := &dataImporter{
		saver:     saver,
		extractor: extractor,
		batchSize: defaultBatchSize,
		logger:    logger,
	}

	for _, o := range opts {
		o(dImporter)
	}

	return dImporter
}

func (i *dataImporter) Import(ctx context.Context) error {
	startTime := time.Now()
	i.logger.Infoln("start importing...")

	extractStartTime := time.Now()
	i.logger.Infoln("start extracting data...")
	data, err := i.extractor.Extract()
	if err != nil {
		return fmt.Errorf("extract data error: %w", err)
	}
	i.logger.Infof("extraction completed successfully, time elapsed: %s", time.Since(extractStartTime))
	i.logger.Infof(formatStatistics(data.Stats))

	saveStartTime := time.Now()
	i.logger.Infof("batch size: %d, start saving data...", i.batchSize)
	if err = i.saver.Save(ctx, data.Rows, i.batchSize); err != nil {
		return fmt.Errorf("save data error: %w", err)
	}
	i.logger.Infof("saving completed successfully, time elapsed: %s", time.Since(saveStartTime))

	i.logger.Infof("import completed successfully, time elapsed: %s", time.Since(startTime))
	return nil
}

func formatStatistics(stats Statics) string {
	return fmt.Sprintf(
		"extracted rows: %d, duplicates: %d, invalid IPs: %d, invalid country codes: %d, invalid counties: %d, invalid citis: %d, invalid latitudes: %d, invalid longitudes: %d, invalid mystery values: %d",
		stats.RawRows, stats.Doubles, stats.BadIPAddress, stats.BadCountryCode, stats.BadCountry, stats.BadCity, stats.BadLatitude, stats.BadLongitude, stats.BadMysteryValue,
	)
}
