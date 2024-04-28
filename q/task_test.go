package q_test

import (
	"testing"

	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/require"
)

func TestNewTas(t *testing.T) {
	cannotMarshall := make(chan int)
	canMarshall := struct{ message string }{message: "hello!"}

	task := q.NewTask("test", canMarshall)
	require.IsType(t, &q.Task{}, task)
	require.Equal(t, "test", task.GetTypename())
	require.Equal(t, canMarshall, task.GetPayload())
	require.Len(t, task.GetOptions(), 0)

	task = q.NewTask("test", cannotMarshall)
	require.IsType(t, &q.Task{}, task)
	require.Equal(t, "test", task.GetTypename())
	require.Equal(t, cannotMarshall, task.GetPayload())
	require.Len(t, task.GetOptions(), 0)
}
