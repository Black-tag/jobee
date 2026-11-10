package task_test

import (
	"testing"

	"jobee/task"
	"jobee/task/memory"
)

func TestRepository_CRUD(t *testing.T) {
	store := memory.New()
	created, err := store.Create(task.Spec{Name: "a"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(created.ID())
	if err != nil || got.Name() != "a" {
		t.Fatalf("get=%v err=%v", got, err)
	}
	edited, err := store.Edit(created.ID(), task.Spec{Name: "b"})
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

func TestRepository_List(t *testing.T) {
	store := memory.New()
	if list, err := store.List(); err != nil || len(list) != 0 {
		t.Fatalf("list=%v err=%v", list, err)
	}
	a, _ := store.Create(task.Spec{Name: "a"})
	b, _ := store.Create(task.Spec{Name: "b"})
	list, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("len=%d", len(list))
	}
	if list[0].ID() != a.ID() || list[1].ID() != b.ID() {
		t.Fatalf("order: %s, %s", list[0].ID(), list[1].ID())
	}
	if err := store.Delete(a.ID()); err != nil {
		t.Fatal(err)
	}
	list, err = store.List()
	if err != nil || len(list) != 1 || list[0].ID() != b.ID() {
		t.Fatalf("after delete list=%v err=%v", list, err)
	}
}

func TestTask_Lifecycle(t *testing.T) {
	store := memory.New()
	tsk, err := store.Create(task.Spec{Name: "step"})
	if err != nil {
		t.Fatal(err)
	}
	if tsk.Status() != task.StatusPending {
		t.Fatalf("status=%s", tsk.Status())
	}
	if err := tsk.Run(); err != nil {
		t.Fatal(err)
	}
	if err := tsk.Pause(); err != nil {
		t.Fatal(err)
	}
	if tsk.Status() != task.StatusPaused {
		t.Fatalf("status=%s", tsk.Status())
	}
	if err := tsk.Run(); err != nil {
		t.Fatal(err)
	}
	if err := tsk.End(); err != nil {
		t.Fatal(err)
	}
	if tsk.Status() != task.StatusEnded {
		t.Fatalf("status=%s", tsk.Status())
	}
	if err := tsk.Run(); err == nil {
		t.Fatal("expected run to fail after end")
	}
}

func TestTask_Fail(t *testing.T) {
	store := memory.New()
	tsk, _ := store.Create(task.Spec{Name: "step"})
	_ = tsk.Run()
	if err := tsk.Fail(); err != nil {
		t.Fatal(err)
	}
	if tsk.Status() != task.StatusError {
		t.Fatalf("status=%s", tsk.Status())
	}
}
