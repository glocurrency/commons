package audit_test

import (
	"testing"

	"github.com/glocurrency/commons/audit"
	"github.com/stretchr/testify/assert"
)

func TestTryMarshall(t *testing.T) {
	cannotMarshall := make(chan int)
	canMarshall := struct{ message string }{message: "hello!"}

	assert.Nil(t, audit.TryMarshall(cannotMarshall))
	assert.NotNil(t, audit.TryMarshall(canMarshall))
}
