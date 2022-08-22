package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/nimusp/geolocation_service/internal/storage"
)

type db struct {
	conn *sql.DB
}

func New() *db {
	return &db{conn: newConn()}
}

func (d *db) Insert(ctx context.Context, data []storage.GeoLocation, batchSize uint) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction err: %w", err)
	}

	defer func() {
		if err != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("rollback err: %w; transaction err: %s", e, err)
			}
			return
		}

		if e := tx.Commit(); e != nil {
			err = fmt.Errorf("commit err: %w", e)
		}
	}()

	chunks := splitToChunks(data, batchSize)
	for _, chunk := range chunks {
		qTemplate := insertTemplate()
		for _, item := range chunk {
			qTemplate = qTemplate.Values(item.IPAddress, item.CountryCode, item.Country, item.City, item.Latitude, item.Longitude, item.MysteryValue)
			query, args, err := qTemplate.ToSql()
			if err != nil {
				return err
			}

			if _, err = tx.ExecContext(ctx, query, args...); err != nil {
				return err
			}
		}
	}

	return nil
}

func splitToChunks(data []storage.GeoLocation, batchSize uint) [][]storage.GeoLocation {
	if len(data) <= int(batchSize) {
		return [][]storage.GeoLocation{data}
	}

	counter := len(data) / int(batchSize)
	if len(data)%int(batchSize) > 0 {
		counter++
	}

	chunks := make([][]storage.GeoLocation, 0, counter)
	for i := 0; i < len(data); i += int(batchSize) {
		top := i + int(batchSize)
		if top > len(data) {
			top = len(data)
		}

		chunks = append(chunks, data[i:top])
	}

	return chunks
}

func insertTemplate() sq.InsertBuilder {
	return sq.Insert("geolocation").Columns(
		"ip_address",
		"country_code",
		"country",
		"city",
		"latitude",
		"longitude",
		"mystery_value").Suffix(`
		ON CONFLICT (ip_address) DO UPDATE SET 
		country_code=EXCLUDED.country_code, 
		country=EXCLUDED.country, 
		city=EXCLUDED.city, 
		latitude=EXCLUDED.latitude, 
		longitude=EXCLUDED.longitude, 
		mystery_value=EXCLUDED.mystery_value`,
	).PlaceholderFormat(sq.Dollar)
}

func (d *db) Select(ctx context.Context, ipAddress string) (*storage.GeoLocation, error) {
	qTemplate := sq.Select(
		"ip_address",
		"country_code",
		"country",
		"city",
		"latitude",
		"longitude",
		"mystery_value",
	).From("geolocation").
		Where(sq.Eq{"ip_address": ipAddress}).PlaceholderFormat(sq.Dollar)

	query, args, err := qTemplate.ToSql()
	if err != nil {
		return nil, err
	}

	var m storage.GeoLocation
	if err = d.conn.QueryRowContext(ctx, query, args...).Scan(&m.IPAddress, &m.CountryCode, &m.Country, &m.City, &m.Latitude, &m.Longitude, &m.MysteryValue); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", err, storage.ErrNoData)
		}

		return nil, err
	}

	return &m, nil
}
