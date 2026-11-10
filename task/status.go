package task

// Status is the lifecycle state of a task.
type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusPaused  Status = "paused"
	StatusEnded   Status = "ended"
	StatusError   Status = "error"
)
