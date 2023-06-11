package logger

import (
	"context"
	"os"

	"github.com/brokeyourbike/nrlogrus"
	"github.com/glocurrency/commons/audit"
	"github.com/sirupsen/logrus"
)

const AuditEventField = "audit_event"

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

func WithContext(ctx context.Context) *Entry {
	entry := &Entry{logger.WithContext(ctx)}
	return entry
}

type Entry struct {
	*logrus.Entry
}

func (e *Entry) WithAuditEvent(event *audit.BasicEvent) *Entry {
	return &Entry{e.WithField(AuditEventField, event)}
}
