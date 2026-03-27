package q_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/assert"
)

func TestTaskOptions(t *testing.T) {
	t.Run("UniqueKey Option", func(t *testing.T) {
		opt := q.UniqueKey("my-unique-id")

		assert.Equal(t, q.UniqueKeyOpt, opt.Type())
		assert.Equal(t, "my-unique-id", opt.Value())
		assert.Equal(t, `UniqueKey("my-unique-id")`, opt.String())
	})

	t.Run("Group Option", func(t *testing.T) {
		opt := q.Group("my-group-id")

		assert.Equal(t, q.GroupOpt, opt.Type())
		assert.Equal(t, "my-group-id", opt.Value())
		assert.Equal(t, `Group("my-group-id")`, opt.String())
	})
}

func TestPubSubOptions(t *testing.T) {
	t.Run("OrderedByTaskName Option", func(t *testing.T) {
		opt := q.OrderedByTaskName()

		assert.Equal(t, q.OrderedByTaskNameOpt, opt.Type())
		assert.Equal(t, true, opt.Value())
		assert.Equal(t, "OrderedByTaskName()", opt.String())
	})

	t.Run("OrderedKey Option", func(t *testing.T) {
		opt := q.OrderedKey("my-order-key")

		assert.Equal(t, q.OrderedKeyOpt, opt.Type())
		assert.Equal(t, "my-order-key", opt.Value())
		assert.Equal(t, `OrderedKey("my-order-key")`, opt.String())
	})
}

func TestCloudTasksOptions(t *testing.T) {
	t.Run("ProcessAt Option", func(t *testing.T) {
		// Use a fixed time to ensure deterministic string formatting
		fixedTime := time.Date(2026, time.March, 27, 20, 0, 0, 0, time.UTC)
		opt := q.ProcessAt(fixedTime)

		assert.Equal(t, q.ProcessAtOpt, opt.Type())
		assert.Equal(t, fixedTime, opt.Value())

		// The string representation relies on time.UnixDate formatting
		expectedStr := fmt.Sprintf("ProcessAt(%v)", fixedTime.Format(time.UnixDate))
		assert.Equal(t, expectedStr, opt.String())
	})
}
