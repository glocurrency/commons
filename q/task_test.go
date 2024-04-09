package q_test

import (
	"testing"

	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/require"
)

func TestNewTaskWithJSON(t *testing.T) {
	cannotMarshall := make(chan int)
	canMarshall := struct{ message string }{message: "hello!"}

	task := q.NewTask("test", canMarshall)
	require.IsType(t, &q.Task{}, task)

	task = q.NewTask("test", cannotMarshall)
	require.Nil(t, task)
}
