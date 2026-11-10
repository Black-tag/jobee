package task

import "errors"

// ApplyRun returns the next status after Run, or an error if the transition is invalid.
func ApplyRun(current Status) (Status, error) {
	switch current {
	case StatusPending, StatusPaused:
		return StatusRunning, nil
	case StatusRunning:
		return "", errors.New("task already running")
	case StatusEnded:
		return "", errors.New("task has ended")
	case StatusError:
		return "", errors.New("task is in error")
	default:
		return "", errors.New("cannot run task")
	}
}

// ApplyPause returns the next status after Pause.
func ApplyPause(current Status) (Status, error) {
	if current != StatusRunning {
		return "", errors.New("task is not running")
	}
	return StatusPaused, nil
}

// ApplyEnd returns the next status after End.
func ApplyEnd(current Status) (Status, error) {
	if current != StatusRunning {
		return "", errors.New("task is not running")
	}
	return StatusEnded, nil
}

// ApplyFail returns the next status after Fail.
func ApplyFail(current Status) (Status, error) {
	switch current {
	case StatusRunning, StatusPaused:
		return StatusError, nil
	case StatusError:
		return "", errors.New("task already in error")
	case StatusEnded:
		return "", errors.New("task has ended")
	default:
		return "", errors.New("cannot fail task")
	}
}
