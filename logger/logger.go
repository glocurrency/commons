package logger

import (
	"context"
	"os"

	"github.com/brokeyourbike/nrlogrus"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(nrlogrus.NewFormatterFromEnvironment(&logrus.JSONFormatter{}))
}

func Log() *logrus.Logger {
	return logger
}

func WithContext(ctx context.Context) *logrus.Entry {
	return logger.WithContext(ctx)
}
