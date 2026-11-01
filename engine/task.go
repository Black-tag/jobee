package engine

import "errors"

// Task is a unit of work that can be attached to a workflow.
type Task interface {
	ID() string
	Name() string
}

// TaskSpec is input for creating or updating a task.
type TaskSpec struct {
	Name string
}

// TaskRepository manages task definitions (create, edit, delete, get).
type TaskRepository interface {
	Create(spec TaskSpec) (Task, error)
	Edit(id string, spec TaskSpec) (Task, error)
	Delete(id string) error
	Get(id string) (Task, error)
}

// basicTask is the default in-memory Task implementation.
type basicTask struct {
	id   string
	name string
}

func (t *basicTask) ID() string   { return t.id }
func (t *basicTask) Name() string { return t.name }

func validateTaskSpec(spec TaskSpec) error {
	if spec.Name == "" {
		return errors.New("task name is required")
	}
	return nil
}
