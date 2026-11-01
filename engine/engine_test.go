package engine

import "testing"

func TestTaskRepository_CRUD(t *testing.T) {
	store := NewMemoryTaskStore()
	created, err := store.Create(TaskSpec{Name: "a"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(created.ID())
	if err != nil || got.Name() != "a" {
		t.Fatalf("get=%v err=%v", got, err)
	}
	edited, err := store.Edit(created.ID(), TaskSpec{Name: "b"})
	if err != nil || edited.Name() != "b" {
		t.Fatalf("edit=%v err=%v", edited, err)
	}
	if err := store.Delete(created.ID()); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(created.ID()); err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestWorkflow_Lifecycle(t *testing.T) {
	store := NewMemoryTaskStore()
	task, _ := store.Create(TaskSpec{Name: "step"})
	wf := NewWorkflow("wf-1", store)

	if err := wf.Add(task.ID()); err != nil {
		t.Fatal(err)
	}
	if err := wf.Start(); err != nil {
		t.Fatal(err)
	}
	if err := wf.Add(task.ID()); err == nil {
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
