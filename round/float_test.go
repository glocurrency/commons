package round_test

import (
	"testing"

	"github.com/glocurrency/commons/round"
	"github.com/stretchr/testify/assert"
)

func TestFloat(t *testing.T) {
	number := 12.3456789

	assert.Equal(t, 12.35, round.Float(number, 2))
	assert.Equal(t, 12.3457, round.Float(number, 4))
}
