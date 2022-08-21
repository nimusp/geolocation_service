package database

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	defaultDriverName     = "postgres"
	defatulMigrationsPath = "./migrations"
)

func newConn() *sqlx.DB {
	connStr := buildConnString()
	conn := sqlx.MustConnect(defaultDriverName, connStr)
	prepareSchema(conn)
	return conn
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

	return sb.String()
}

func prepareSchema(conn *sqlx.DB) {
	queryes, err := readMigrations(defatulMigrationsPath)
	if err != nil {
		panic(err)
	}

	for _, query := range queryes {
		_ = conn.MustExec(query)
	}
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
