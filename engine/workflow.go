package engine

import (
	"errors"

	"jobee/task"
)

// Workflow owns an ordered list of tasks and exposes lifecycle controls.
type Workflow interface {
	ID() string
	Status() WorkflowStatus
	TaskIDs() []string

	Start() error
	Add(taskID string) error
	Cancel() error
	Pause() error
}

// BasicWorkflow is an in-process Workflow implementation (task order not persisted yet).
type BasicWorkflow struct {
	id     string
	status WorkflowStatus
	order  []string
	tasks  task.Repository
}

func NewWorkflow(id string, tasks task.Repository) *BasicWorkflow {
	return &BasicWorkflow{
		id:     id,
		status: WorkflowDraft,
		tasks:  tasks,
	}
}

func (w *BasicWorkflow) ID() string { return w.id }

func (w *BasicWorkflow) Status() WorkflowStatus { return w.status }

func (w *BasicWorkflow) TaskIDs() []string {
	out := make([]string, len(w.order))
	copy(out, w.order)
	return out
}

func (w *BasicWorkflow) Start() error {
	switch w.status {
	case WorkflowDraft, WorkflowPaused:
		w.status = WorkflowRunning
		return nil
	case WorkflowRunning:
		return errors.New("workflow already running")
	case WorkflowCancelled:
		return errors.New("workflow is cancelled")
	default:
		return errors.New("cannot start workflow")
	}
}

func (w *BasicWorkflow) Add(taskID string) error {
	if w.status == WorkflowRunning {
		return errors.New("cannot add tasks while workflow is running")
	}
	if w.status == WorkflowCancelled {
		return errors.New("workflow is cancelled")
	}
	if _, err := w.tasks.Get(taskID); err != nil {
		return err
	}
	for _, id := range w.order {
		if id == taskID {
			return errors.New("task already in workflow")
		}
	}
	w.order = append(w.order, taskID)
	return nil
}

func (w *BasicWorkflow) Cancel() error {
	switch w.status {
	case WorkflowRunning, WorkflowPaused, WorkflowDraft:
		w.status = WorkflowCancelled
		return nil
	default:
		return errors.New("workflow already cancelled")
	}
}

func (w *BasicWorkflow) Pause() error {
	if w.status != WorkflowRunning {
		return errors.New("workflow is not running")
	}
	w.status = WorkflowPaused
	return nil
}
