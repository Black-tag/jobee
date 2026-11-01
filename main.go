package main

import (
	"fmt"
	"log"

	"jobee/engine"
)

func main() {
	tasks := engine.NewMemoryTaskStore()

	validate, err := tasks.Create(engine.TaskSpec{Name: "validate"})
	if err != nil {
		log.Fatal(err)
	}
	ship, err := tasks.Create(engine.TaskSpec{Name: "ship"})
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

	edited, err := tasks.Edit(validate.ID(), engine.TaskSpec{Name: "validate-order"})
	if err != nil {
		log.Fatal(err)
	}
	got, _ := tasks.Get(edited.ID())
	fmt.Println("task get:", got.ID(), got.Name())
}
