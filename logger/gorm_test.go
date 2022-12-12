package logger_test

import (
	"testing"
	"time"

	"github.com/glocurrency/commons/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewGormLogger(t *testing.T) {
	log := logger.Log()

	g := logger.NewGormLogger(log, 2*time.Second, true)
	assert.Same(t, log, g.Logger)
	assert.Equal(t, 2*time.Second, g.SlowThreshold)
	assert.Equal(t, true, g.SkipErrRecordNotFound)
}
