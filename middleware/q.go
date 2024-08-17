package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/instrumentation"
	"github.com/glocurrency/commons/q"
)

const QMessageCtxKey = "qMessageCtx"

type qctx struct {
	locker Locker
}

func NewQCtx(locker Locker) *qctx {
	return &qctx{locker: locker}
}

func (m *qctx) RequireValidMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		msg, err := q.NewQMessage(ctx.Request)
		if err != nil {
			instrumentation.NoticeError(ctx, err, "cannot parse queue message")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		instrumentation.NoticeInfo(ctx, "message from queue received",
			instrumentation.WithField("queue_message", msg))

		if msg.UniqueKey != "" {
			if err := m.locker.TryToLock(ctx, msg.UniqueKey); err != nil {
				instrumentation.NoticeError(ctx, err, "cannot lock task",
					instrumentation.WithField("unique_key", msg.UniqueKey))
				ctx.AbortWithStatus(http.StatusUnprocessableEntity)
				return
			}
		}

		ctx.Set(QMessageCtxKey, msg)
		ctx.Next()
	}
}

// MustGetQMessageFromContext returns the Q message from the context.
func MustGetQMessageFromContext(ctx *gin.Context) q.QMessage {
	msg := ctx.MustGet(QMessageCtxKey).(q.QMessage)
	return msg
}
