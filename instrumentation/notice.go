package instrumentation

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/glocurrency/commons/logger"
)

func NoticeError(ctx context.Context, err error, msg string, opts ...NoticeOption) {
	entry := ApplyOptions(logger.WithContext(ctx), opts...)
	entry.WithError(err).Error(msg)

	if hub := getHubFromContext(ctx); hub != nil {
		hub.CaptureException(err)
		return
	}

	sentry.CaptureException(err)
}

func NoticeWarning(ctx context.Context, msg string, opts ...NoticeOption) {
	entry := ApplyOptions(logger.WithContext(ctx), opts...)
	entry.Warn(msg)
}

func NoticeInfo(ctx context.Context, msg string, opts ...NoticeOption) {
	entry := ApplyOptions(logger.WithContext(ctx), opts...)
	entry.Info(msg)
}
