package task

// Repository manages tasks (create, edit, run, delete, get, list).
type Repository interface {
	Create(spec Spec) (Task, error)
	Edit(id string, spec Spec) (Task, error)
	Run(id string) (Task, error)
	Delete(id string) error
	Get(id string) (Task, error)
	List() ([]Task, error)
}
