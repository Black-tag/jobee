package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"jobee/engine"
	"jobee/internal/db"
	"jobee/task"
	"jobee/task/memory"
	"jobee/task/postgres"
)

func main() {
	tasks, cleanup, err := openTaskRepository()
	if err != nil {
		log.Fatal(err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	validate, err := tasks.Create(task.Spec{Name: "validate"})
	if err != nil {
		log.Fatal(err)
	}
	ship, err := tasks.Create(task.Spec{Name: "ship"})
	if err != nil {
		log.Fatal(err)
	}

	wf := engine.NewWorkflow("order-flow", tasks)
	var workflow engine.Workflow = wf

	if err := workflow.Add(validate.ID()); err != nil {
		log.Fatal(err)
	}
	if err := workflow.Add(ship.ID()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("workflow", workflow.ID(), "status", workflow.Status(), "tasks", workflow.TaskIDs())

	if err := workflow.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("after start:", workflow.Status())

	if err := workflow.Pause(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("after pause:", workflow.Status())

	if err := workflow.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("after resume:", workflow.Status())

	edited, err := tasks.Edit(validate.ID(), task.Spec{Name: "validate-order"})
	if err != nil {
		log.Fatal(err)
	}
	got, _ := tasks.Get(edited.ID())
	fmt.Println("task get:", got.ID(), got.Name())
}

func openTaskRepository() (task.Repository, func(), error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return memory.New(), nil, nil
	}
	ctx := context.Background()
	if err := db.Up(dsn); err != nil {
		return nil, nil, err
	}
	pool, err := db.OpenPool(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}
	return postgres.New(pool), func() { pool.Close() }, nil
}
