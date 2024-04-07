package logger

import (
	"context"

	"github.com/glocurrency/commons/audit"
	"github.com/sirupsen/logrus"
)

const AuditEventField = "audit_event"

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetOutput(Writer{})
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
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

func (e *Entry) EWithFields(fields map[string]interface{}) *Entry {
	return &Entry{e.WithFields(fields)}
}

func (e *Entry) WithAuditEvent(event *audit.BasicEvent) *Entry {
	return &Entry{e.WithField(AuditEventField, event)}
}

func (e *Entry) PushWithAuditEvent(event *audit.BasicEvent) {
	Entry{e.WithField(AuditEventField, event)}.Info(AuditEventField)
}
