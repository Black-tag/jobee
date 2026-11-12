package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"jobee/internal/db/sqlc"
	"jobee/task"
)

// Store is a PostgreSQL-backed task.Repository (sqlc + pgx).
type Store struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool, q: sqlc.New(pool)}
}

func (s *Store) Create(spec task.Spec) (task.Task, error) {
	if err := task.ValidateSpec(spec); err != nil {
		return nil, err
	}
	row, err := s.q.CreateTask(context.Background(), spec.Name)
	if err != nil {
		return nil, err
	}
	return s.rowToTask(row)
}

func (s *Store) Edit(id string, spec task.Spec) (task.Task, error) {
	if err := task.ValidateSpec(spec); err != nil {
		return nil, err
	}
	pid, err := parseID(id)
	if err != nil {
		return nil, err
	}
	row, err := s.q.UpdateTaskName(context.Background(), sqlc.UpdateTaskNameParams{
		ID:   pid,
		Name: spec.Name,
	})
	if err != nil {
		return nil, mapNotFound(err)
	}
	return s.rowToTask(row)
}

func (s *Store) Delete(id string) error {
	pid, err := parseID(id)
	if err != nil {
		return err
	}
	if _, err := s.q.GetTask(context.Background(), pid); err != nil {
		return mapNotFound(err)
	}
	return s.q.DeleteTask(context.Background(), pid)
}

func (s *Store) Run(id string) (task.Task, error) {
	t, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	if err := t.Run(); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Store) Get(id string) (task.Task, error) {
	pid, err := parseID(id)
	if err != nil {
		return nil, err
	}
	row, err := s.q.GetTask(context.Background(), pid)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return s.rowToTask(row)
}

func (s *Store) List() ([]task.Task, error) {
	rows, err := s.q.ListTasks(context.Background())
	if err != nil {
		return nil, err
	}
	out := make([]task.Task, 0, len(rows))
	for _, row := range rows {
		t, err := s.rowToTask(row)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (s *Store) updateStatus(pid pgtype.UUID, status task.Status) (sqlc.Task, error) {
	return s.q.UpdateTaskStatus(context.Background(), sqlc.UpdateTaskStatusParams{
		ID:     pid,
		Status: string(status),
	})
}

func (s *Store) rowToTask(row sqlc.Task) (*rowTask, error) {
	if !row.ID.Valid {
		return nil, fmt.Errorf("task row missing id")
	}
	created, updated := rowTimestamps(row)
	return &rowTask{
		store:     s,
		id:        row.ID,
		idStr:     pgUUIDString(row.ID),
		name:      row.Name,
		status:    task.Status(row.Status),
		createdAt: created,
		updatedAt: updated,
	}, nil
}

type rowTask struct {
	store     *Store
	id        pgtype.UUID
	idStr     string
	name      string
	status    task.Status
	createdAt time.Time
	updatedAt time.Time
}

func (t *rowTask) ID() string            { return t.idStr }
func (t *rowTask) Name() string          { return t.name }
func (t *rowTask) Status() task.Status   { return t.status }
func (t *rowTask) CreatedAt() time.Time  { return t.createdAt }
func (t *rowTask) UpdatedAt() time.Time  { return t.updatedAt }

func (t *rowTask) Run() error    { return t.transition(task.ApplyRun) }
func (t *rowTask) Pause() error  { return t.transition(task.ApplyPause) }
func (t *rowTask) End() error    { return t.transition(task.ApplyEnd) }
func (t *rowTask) Fail() error   { return t.transition(task.ApplyFail) }

func (t *rowTask) transition(apply func(task.Status) (task.Status, error)) error {
	next, err := apply(t.status)
	if err != nil {
		return err
	}
	row, err := t.store.updateStatus(t.id, next)
	if err != nil {
		return err
	}
	t.name = row.Name
	t.status = task.Status(row.Status)
	t.createdAt, t.updatedAt = rowTimestamps(row)
	return nil
}

func rowTimestamps(row sqlc.Task) (time.Time, time.Time) {
	var created, updated time.Time
	if row.CreatedAt.Valid {
		created = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		updated = row.UpdatedAt.Time
	}
	return created, updated
}

func parseID(id string) (pgtype.UUID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, errors.New("task not found")
	}
	return pgtype.UUID{Bytes: u, Valid: true}, nil
}

func pgUUIDString(p pgtype.UUID) string {
	if !p.Valid {
		return ""
	}
	u, err := uuid.FromBytes(p.Bytes[:])
	if err != nil {
		return ""
	}
	return u.String()
}

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("task not found")
	}
	return err
}
