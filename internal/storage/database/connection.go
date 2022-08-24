package database

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	defaultDriverName     = "postgres"
	defatulMigrationsPath = "/internal/storage/database/migrations"

	retryNumber = 5
	retryPause  = time.Second
)

func newConn() (*sql.DB, error) {
	connStr := buildConnString()
	conn, err := sql.Open(defaultDriverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("open connection error: %w", err)
	}

	for i := 0; i < retryNumber; i++ {
		if err = conn.Ping(); err == nil {
			break
		}
		time.Sleep(retryPause)
	}
	if err != nil {
		return nil, fmt.Errorf("ping DB error: %w", err)
	}

	if err = prepareSchema(conn); err != nil {
		return nil, fmt.Errorf("prepare DB schema error: %w", err)
	}

	return conn, nil
}

func buildConnString() string {
	var sb strings.Builder
	sb.WriteString("sslmode=disable")
	if host := os.Getenv(envHost); len(host) > 0 {
		sb.WriteString(" host=")
		sb.WriteString(host)
	}

	if port := os.Getenv(envPort); len(port) > 0 {
		sb.WriteString(" port=")
		sb.WriteString(port)
	}

	if user := os.Getenv(envUser); len(user) > 0 {
		sb.WriteString(" user=")
		sb.WriteString(user)
	}

	if password := os.Getenv(envPassword); len(password) > 0 {
		sb.WriteString(" password=")
		sb.WriteString(password)
	}

	if dbName := os.Getenv(envDbName); len(dbName) > 0 {
		sb.WriteString(" dbname=")
		sb.WriteString(dbName)
	}

	return sb.String()
}

func prepareSchema(conn *sql.DB) error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("can't get dir path: %w", err)
	}

	queryes, err := readMigrations(filepath.Join(dir, defatulMigrationsPath))
	if err != nil {
		panic(err)
	}

	for _, query := range queryes {
		if _, err := conn.Exec(query); err != nil {
			return fmt.Errorf("execute query (%s) error: %w", query, err)
		}
	}

	return nil
}

func readMigrations(migraionPath string) ([]string, error) {
	entries, err := os.ReadDir(migraionPath)
	if err != nil {
		return nil, fmt.Errorf("open migration dir err: %w", err)
	}

	if len(entries) == 0 {
		return []string{}, nil
	}

	content := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		path := filepath.Join(migraionPath, fileName)
		rawData, err := readFile(path)
		if err != nil {
			return nil, err
		}

		content = append(content, string(rawData))
	}

	return content, nil
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s error: %w", path, err)
	}
	defer f.Close()

	return io.ReadAll(f)
}
