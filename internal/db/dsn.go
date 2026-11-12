package db

import "os"

// DSN returns PostgreSQL connection string from DATABASE_URL or a local default.
func DSN() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://postgres:postgres@127.0.0.1:5432/jobee?sslmode=disable"
}
