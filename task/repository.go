package task

// Repository manages tasks (create, edit, delete, get, list).
type Repository interface {
	Create(spec Spec) (Task, error)
	Edit(id string, spec Spec) (Task, error)
	Delete(id string) error
	Get(id string) (Task, error)
	List() ([]Task, error)
}
