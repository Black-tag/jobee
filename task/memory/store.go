package memory

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"jobee/task"
)

// Store is an in-memory task.Repository.
type Store struct {
	mu     sync.RWMutex
	nextID int
	tasks  map[string]*record
}

func New() *Store {
	return &Store{tasks: make(map[string]*record)}
}

type record struct {
	id     string
	name   string
	status task.Status
}

func (r *record) ID() string          { return r.id }
func (r *record) Name() string        { return r.name }
func (r *record) Status() task.Status { return r.status }

func (r *record) Run() error {
	next, err := task.ApplyRun(r.status)
	if err != nil {
		return err
	}
	r.status = next
	return nil
}

func (r *record) Pause() error {
	next, err := task.ApplyPause(r.status)
	if err != nil {
		return err
	}
	r.status = next
	return nil
}

func (r *record) End() error {
	next, err := task.ApplyEnd(r.status)
	if err != nil {
		return err
	}
	r.status = next
	return nil
}

func (r *record) Fail() error {
	next, err := task.ApplyFail(r.status)
	if err != nil {
		return err
	}
	r.status = next
	return nil
}

func (s *Store) Create(spec task.Spec) (task.Task, error) {
	if err := task.ValidateSpec(spec); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	id := fmt.Sprintf("task-%d", s.nextID)
	r := &record{id: id, name: spec.Name, status: task.StatusPending}
	s.tasks[id] = r
	return r, nil
}

func (s *Store) Edit(id string, spec task.Spec) (task.Task, error) {
	if err := task.ValidateSpec(spec); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	r.name = spec.Name
	return r, nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return errors.New("task not found")
	}
	delete(s.tasks, id)
	return nil
}

func (s *Store) Get(id string) (task.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	return r, nil
}

func (s *Store) List() ([]task.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]string, 0, len(s.tasks))
	for id := range s.tasks {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]task.Task, len(ids))
	for i, id := range ids {
		out[i] = s.tasks[id]
	}
	return out, nil
}
