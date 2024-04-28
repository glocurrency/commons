package q

// Task represents a unit of work to be performed.
type Task struct {
	// typename indicates the type of task to be performed.
	typename string

	// payload holds data needed to perform the task.
	payload interface{}

	// opts holds options for the task.
	opts []TaskOption
}

func (t Task) GetTypename() string {
	return t.typename
}

func (t Task) GetPayload() interface{} {
	return t.payload
}

func (t Task) GetOptions() []TaskOption {
	return t.opts
}

// NewTask returns a new Task given a type name and payload data that will be marshaled.
// Options can be passed to configure task processing behavior.
func NewTask(typename string, payload interface{}, opts ...TaskOption) *Task {
	return &Task{
		typename: typename,
		payload:  payload,
		opts:     opts,
	}
}

type TaskInfo struct {
	ID string
}
