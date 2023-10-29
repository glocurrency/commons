package logger_test

import (
	"context"
	"testing"

	"github.com/glocurrency/commons/audit"
	"github.com/glocurrency/commons/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

type mockTarget struct{ ID uuid.UUID }

func (m mockTarget) GetID() uuid.UUID {
	return m.ID
}

func (m mockTarget) GetAuditTargetType() audit.TargetType {
	return "mock-target"
}

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
	target := mockTarget{ID: uuid.New()}

	event := audit.NewBasicEvent(
		"audit-type",
		target,
		"actor-type",
		audit.WithPayload(target),
	)
	assert.NotNil(t, event)

	hook := test.NewLocal(logger.Log())

	logger.WithContext(context.Background()).WithAuditEvent(event).Info("test")
	assert.Len(t, hook.Entries, 1)
}
