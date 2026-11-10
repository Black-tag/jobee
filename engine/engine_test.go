package engine

import (
	"testing"

	"jobee/task"
	"jobee/task/memory"
)

func TestWorkflow_Lifecycle(t *testing.T) {
	store := memory.New()
	tsk, _ := store.Create(task.Spec{Name: "step"})
	wf := NewWorkflow("wf-1", store)

	if err := wf.Add(tsk.ID()); err != nil {
		t.Fatal(err)
	}
	if err := wf.Start(); err != nil {
		t.Fatal(err)
	}
	if err := wf.Add(tsk.ID()); err == nil {
		t.Fatal("expected add to fail while running")
	}
	if err := wf.Pause(); err != nil {
		t.Fatal(err)
	}
	if wf.Status() != WorkflowPaused {
		t.Fatalf("status=%s", wf.Status())
	}
	if err := wf.Cancel(); err != nil {
		t.Fatal(err)
	}
}
