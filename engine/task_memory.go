package engine

import (
	"errors"
	"fmt"
	"sync"
)

// MemoryTaskStore is an in-memory TaskRepository.
type MemoryTaskStore struct {
	mu     sync.RWMutex
	nextID int
	tasks  map[string]*basicTask
}

func NewMemoryTaskStore() *MemoryTaskStore {
	return &MemoryTaskStore{tasks: make(map[string]*basicTask)}
}

func (s *MemoryTaskStore) Create(spec TaskSpec) (Task, error) {
	if err := validateTaskSpec(spec); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	id := fmt.Sprintf("task-%d", s.nextID)
	t := &basicTask{id: id, name: spec.Name}
	s.tasks[id] = t
	return t, nil
}

func (s *MemoryTaskStore) Edit(id string, spec TaskSpec) (Task, error) {
	if err := validateTaskSpec(spec); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	t.name = spec.Name
	return t, nil
}

func (s *MemoryTaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return errors.New("task not found")
	}
	delete(s.tasks, id)
	return nil
}

func (s *MemoryTaskStore) Get(id string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	return t, nil
}
