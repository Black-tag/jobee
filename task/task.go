package task

import "errors"

// Task is a unit of work that can be attached to a workflow.
type Task interface {
	ID() string
	Name() string
	Status() Status

	Run() error
	Pause() error
	End() error
	Fail() error
}

// Spec is input for creating or updating a task.
type Spec struct {
	Name string
}

func ValidateSpec(spec Spec) error {
	if spec.Name == "" {
		return errors.New("task name is required")
	}
	return nil
}
