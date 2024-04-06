package instrumentation

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/glocurrency/commons/logger"
)

func NoticeError(ctx context.Context, err error, msg string, opts ...NoticeOption) {
	entry := logger.WithContext(ctx)
	for _, o := range opts {
		o.Apply(entry)
	}
	entry.WithError(err).Error(msg)

	sentry.CaptureException(err)
}

func NoticeWarning(ctx context.Context, msg string, opts ...NoticeOption) {
	entry := logger.WithContext(ctx)
	for _, o := range opts {
		entry = o.Apply(entry)
	}
	entry.Warn(msg)
}

func NoticeInfo(ctx context.Context, msg string, opts ...NoticeOption) {
	entry := logger.WithContext(ctx)
	for _, o := range opts {
		entry = o.Apply(entry)
	}
	entry.Info(msg)
}
