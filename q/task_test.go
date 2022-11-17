package q_test

import (
	"testing"

	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/assert"
)

func TestNewTaskWithJSON(t *testing.T) {
	cannotMarshall := make(chan int)
	canMarshall := struct{ message string }{message: "hello!"}

	task, err := q.NewTaskWithJSON("test", canMarshall)
	assert.NoError(t, err)
	assert.IsType(t, &q.Task{}, task)

	task, err = q.NewTaskWithJSON("test", cannotMarshall)
	assert.Error(t, err)
	assert.Nil(t, task)
}
