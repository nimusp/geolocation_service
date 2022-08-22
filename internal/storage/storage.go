package storage

import (
	"context"
	"errors"

	"github.com/nimusp/geolocation_service/internal/dao"
	"github.com/nimusp/geolocation_service/internal/importer"
)

type storageWrapper struct {
	db Database
}

func NewStorage(db Database) *storageWrapper {
	return &storageWrapper{
		db: db,
	}
}

func (s *storageWrapper) Save(ctx context.Context, data []importer.GeoLocation, batchSize uint) error {
	storageModels := make([]GeoLocation, 0, len(data))
	for _, model := range data {
		storageModels = append(storageModels, GeoLocation{
			IPAddress:    model.IPAddress,
			CountryCode:  model.CountryCode,
			Country:      model.Country,
			City:         model.City,
			Latitude:     model.Latitude,
			Longitude:    model.Longitude,
			MysteryValue: model.MysteryValue,
		})
	}
	return s.db.Insert(ctx, storageModels, batchSize)
}

func (s *storageWrapper) GetByIP(ctx context.Context, ipAddress string) (*dao.GeoLocation, error) {
	model, err := s.db.Select(ctx, ipAddress)
	if err != nil {
		if errors.Is(err, ErrNoData) {
			return nil, dao.ErrNoData
		}
		return nil, err
	}

	return &dao.GeoLocation{
		IPAddress:    model.IPAddress,
		CountryCode:  model.CountryCode,
		Country:      model.Country,
		City:         model.City,
		Latitude:     model.Latitude,
		Longitude:    model.Longitude,
		MysteryValue: model.MysteryValue,
	}, nil
}
