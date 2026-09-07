// Package db manages the Postgres connection and schema migrations.
package db

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed all:migrations
var migrationsFS embed.FS

// Open connects to Postgres and applies any pending migrations.
func Open(databaseURL string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("set migration dialect: %w", err)
	}

	if err := goose.Up(conn, "migrations"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return conn, nil
}
