package logger_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/glocurrency/commons/logger"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	glog "gorm.io/gorm/logger"
)

func TestNewGormLogger(t *testing.T) {
	ctx := context.TODO()
	log := logger.Log()
	hook := logrustest.NewLocal(log)

	g := logger.NewGormLogger(log, 2*time.Second, true)
	require.Same(t, log, g.Logger)
	require.Equal(t, 2*time.Second, g.SlowThreshold)
	require.Equal(t, true, g.SkipErrRecordNotFound)

	require.Equal(t, g, g.LogMode(glog.Error))
	g.Info(ctx, "info")
	g.Warn(ctx, "warn")
	g.Error(ctx, "error")
	g.Trace(ctx, time.Now().Add(-1*time.Minute), func() (string, int64) {
		return "", 0
	}, nil)
	g.Trace(ctx, time.Now().Add(-1*time.Minute), func() (string, int64) {
		return "", 0
	}, errors.New("fail"))

	assert.Equal(t, 5, len(hook.Entries))
}
