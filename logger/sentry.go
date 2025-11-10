package logger

import (
	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/attribute"
	"github.com/sirupsen/logrus"
)

type SentryHook struct{}

func (h SentryHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (h SentryHook) Fire(entry *logrus.Entry) error {
	if entry == nil {
		return nil
	}

	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	sentryLogger := sentry.NewLogger(ctx)

	attributes := make([]attribute.Builder, 0)

	if entry.Data != nil {
		for k, v := range entry.Data {
			switch typed := v.(type) {
			case string:
				attributes = append(attributes, attribute.String(k, typed))
			case int:
				attributes = append(attributes, attribute.Int(k, typed))
			case int64:
				attributes = append(attributes, attribute.Int64(k, typed))
			case float64:
				attributes = append(attributes, attribute.Float64(k, typed))
			case bool:
				attributes = append(attributes, attribute.Bool(k, typed))
			}
		}
	}

	sentryLogger.SetAttributes(attributes...)

	switch entry.Level {
	case logrus.TraceLevel:
		sentryLogger.Trace()
	case logrus.DebugLevel:
		sentryLogger.Debug()
	case logrus.InfoLevel:
		sentryLogger.Info()
	case logrus.WarnLevel:
		sentryLogger.Warn()
	case logrus.ErrorLevel:
		sentryLogger.Error()
	case logrus.FatalLevel:
		sentryLogger.Fatal()
	case logrus.PanicLevel:
		sentryLogger.Panic()
	}

	return nil
}
