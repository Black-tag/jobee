package engine

// WorkflowStatus is the lifecycle state of a workflow definition/run.
type WorkflowStatus string

const (
	WorkflowDraft     WorkflowStatus = "draft"
	WorkflowRunning   WorkflowStatus = "running"
	WorkflowPaused    WorkflowStatus = "paused"
	WorkflowCancelled WorkflowStatus = "cancelled"
)
