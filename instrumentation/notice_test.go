package instrumentation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/glocurrency/commons/instrumentation"
	"github.com/glocurrency/commons/logger"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNoticeInfo(t *testing.T) {
	hook := test.NewLocal(logger.Log())

	instrumentation.NoticeInfo(context.TODO(), "hi!")
	instrumentation.NoticeInfo(context.TODO(), "hi!", instrumentation.WithField("name", "john"))
	instrumentation.NoticeInfo(context.TODO(), "hi!", instrumentation.WithFields(map[string]interface{}{"age": 100}))

	assert.Len(t, hook.Entries, 3)
}

func TestNoticeWarning(t *testing.T) {
	hook := test.NewLocal(logger.Log())

	instrumentation.NoticeWarning(context.TODO(), "hi!")
	instrumentation.NoticeWarning(context.TODO(), "hi!", instrumentation.WithField("name", "john"))
	instrumentation.NoticeWarning(context.TODO(), "hi!", instrumentation.WithFields(map[string]interface{}{"age": 100}))

	assert.Len(t, hook.Entries, 3)
}

func TestNoticeError(t *testing.T) {
	hook := test.NewLocal(logger.Log())

	err := errors.New("i am an error!")

	instrumentation.NoticeError(context.TODO(), err, "error!")
	instrumentation.NoticeError(context.TODO(), err, "error!", instrumentation.WithField("name", "john"))
	instrumentation.NoticeError(context.TODO(), err, "error!", instrumentation.WithFields(map[string]interface{}{"age": 100}))

	assert.Len(t, hook.Entries, 3)
}
