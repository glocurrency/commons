package logger

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Inspired by: https://github.com/onrik/gorm-logrus/logger.go
// MIT License: https://github.com/onrik/gorm-logrus/LICENSE
type GormLogger struct {
	Logger                *logrus.Logger
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
}

func NewGormLogger(logger *logrus.Logger, slowTreshold time.Duration, skipErrNotFound bool) *GormLogger {
	return &GormLogger{
		Logger:                logger,
		SlowThreshold:         slowTreshold,
		SkipErrRecordNotFound: skipErrNotFound,
	}
}

func (l *GormLogger) LogMode(glog.LogLevel) glog.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Infof(s, args...)
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Warnf(s, args...)
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.WithContext(ctx).Errorf(s, args...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	fields := logrus.Fields{"file_with_line_num": utils.FileWithLineNum()}

	sql, rows := fc()
	if rows == -1 {
		fields["rows_affected"] = "-"
	} else {
		fields["rows_affected"] = rows
	}

	elapsed := time.Since(begin)
	fields["elapsed"] = elapsed

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.SkipErrRecordNotFound):
		l.Logger.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold:
		l.Logger.WithContext(ctx).WithFields(fields).Warnf("SLOW SQL >= %v, %s [%s]", l.SlowThreshold, sql, elapsed)
	}
}
