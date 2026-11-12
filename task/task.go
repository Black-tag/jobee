package task

import (
	"errors"
	"time"
)

// Task is a persisted unit of work (catalog row + runtime status).
type Task interface {
	ID() string
	Name() string
	Status() Status
	CreatedAt() time.Time
	UpdatedAt() time.Time

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
