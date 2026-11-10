package postgres_test

import (
	"context"
	"os"
	"testing"

	"jobee/internal/db"
	"jobee/task"
	"jobee/task/postgres"
)

func TestPostgresStore_Integration(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set")
	}
	if err := db.Up(dsn); err != nil {
		t.Fatal(err)
	}
	pool, err := db.OpenPool(context.Background(), dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	store := postgres.New(pool)
	created, err := store.Create(task.Spec{Name: "db-step"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(created.ID())
	if err != nil || got.Name() != "db-step" {
		t.Fatalf("get err=%v name=%q", err, got.Name())
	}
	if err := got.Run(); err != nil {
		t.Fatal(err)
	}
	list, err := store.List()
	if err != nil || len(list) < 1 {
		t.Fatalf("list=%v err=%v", list, err)
	}
	if err := store.Delete(created.ID()); err != nil {
		t.Fatal(err)
	}
}
