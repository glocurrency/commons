package instrumentation

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/glocurrency/commons/logger"
)

// TODO: depending on ctx, insert GCP related info
// sample: https://github.com/ncruces/go-gcp/blob/master/glog/glog.go
// sample: https://github.com/googleapis/google-cloud-go/blob/main/logging/logging.go
// fields: https://cloud.google.com/logging/docs/structured-logging#special-payload-fields

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
