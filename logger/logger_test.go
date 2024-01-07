package logger_test

import (
	"context"
	"testing"

	"github.com/glocurrency/commons/audit"
	"github.com/glocurrency/commons/logger"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReturnSameInstanceOfLog(t *testing.T) {
	logger1 := logger.Log()
	logger2 := logger.Log()

	assert.Same(t, logger1, logger2)
}

func TestWithContext(t *testing.T) {
	ctx := context.TODO()

	e := logger.WithContext(ctx)
	assert.IsType(t, &logger.Entry{}, e)
	assert.Equal(t, ctx, e.Context)
}

func TestLogEntry_WithAuditEvent(t *testing.T) {
	event := audit.NewBasicEvent(
		"audit-type",
		"target-type",
		"actor-type",
	)
	require.NotNil(t, event)

	hook := test.NewLocal(logger.Log())

	logger.WithContext(context.TODO()).WithAuditEvent(event).Info("test")
	logger.WithContext(context.TODO()).PushWithAuditEvent(event)

	assert.Len(t, hook.Entries, 2)
}
