package db

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Up applies all pending database migrations.
func Up(dsn string) error {
	goose.SetBaseFS(migrationFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}
	defer sqlDB.Close()
	return goose.Up(sqlDB, "migrations")
}
