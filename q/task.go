package q

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

type TaskInfo struct {
	ID string
}
