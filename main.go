package main

import (
	"context"
	"fmt"
	"log"

	"jobee/internal/db"
	"jobee/task"
	"jobee/task/postgres"
)

func main() {
	ctx := context.Background()
	dsn := db.DSN()

	if err := db.Up(dsn); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	pool, err := db.OpenPool(ctx, dsn)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()

	repo := postgres.New(pool)

	created, err := repo.Create(task.Spec{Name: "validate-order"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("create:", created.ID(), created.Name(), created.Status())

	edited, err := repo.Edit(created.ID(), task.Spec{Name: "validate"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("edit:", edited.Name())

	running, err := repo.Run(created.ID())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("run:", running.Status())

	got, err := repo.Get(created.ID())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("get:", got.ID(), got.Name(), got.Status())

	if err := repo.Delete(created.ID()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("delete: ok")

	list, err := repo.List()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("list: %d tasks\n", len(list))
}
