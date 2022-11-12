package logger_test

import (
	"context"
	"testing"

	"github.com/glocurrency/commons/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestReturnSameInstanceOfLog(t *testing.T) {
	logger1 := logger.Log()
	logger2 := logger.Log()

	assert.Same(t, logger1, logger2)
}

func TestWithContext(t *testing.T) {
	ctx := context.TODO()

	e := logger.WithContext(ctx)
	assert.IsType(t, &logrus.Entry{}, e)
	assert.Same(t, ctx, e.Context)
}
