package q

import (
	"encoding/json"
	"fmt"
)

// Task represents a unit of work to be performed.
type Task struct {
	// typename indicates the type of task to be performed.
	typename string

	// payload holds data needed to perform the task.
	payload []byte

	// opts holds options for the task.
	opts []TaskOption
}

// NewTask returns a new Task given a type name and payload data.
// Options can be passed to configure task processing behavior.
func NewTask(typename string, payload []byte, opts ...TaskOption) *Task {
	return &Task{
		typename: typename,
		payload:  payload,
		opts:     opts,
	}
}

// NewTaskWithJSON returns a new Task given a type name and JSON payload that will be marshaled.
// Options can be passed to configure task processing behavior.
func NewTaskWithJSON(typename string, payload interface{}, opts ...TaskOption) (*Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return NewTask(typename, data, opts...), nil
}

type TaskInfo struct {
	ID string
}
