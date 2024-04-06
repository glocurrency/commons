package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/instrumentation"
	"github.com/glocurrency/commons/q"
)

const PubSubMessageCtxKey = "pubSubMessageCtx"

type Locker interface {
	TryToLock(ctx context.Context, key string) error
}

type pubSubCtx struct {
	locker Locker
}

func NewPubSubCtx(l Locker) *pubSubCtx {
	return &pubSubCtx{locker: l}
}

func (m *pubSubCtx) RequireValidMessage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var msg q.PubSubMessage
		if err := ctx.ShouldBindJSON(&msg); err != nil {
			instrumentation.NoticeError(ctx, err, "cannot unmarshal pubsub message")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		instrumentation.NoticeInfo(ctx, "pubsub message received",
			instrumentation.WithField("pubsub_message", msg))

		if msg.GetUniqueKey() != "" {
			if err := m.locker.TryToLock(ctx, msg.GetUniqueKey()); err != nil {
				instrumentation.NoticeError(ctx, err, "cannot lock task",
					instrumentation.WithField("unique_key", msg.GetUniqueKey()))
				ctx.AbortWithStatus(http.StatusUnprocessableEntity)
				return
			}
		}

		ctx.Set(PubSubMessageCtxKey, msg)
		ctx.Next()
	}
}

// MustGetMessageFromContext returns the PubSub message from the context.
func MustGetMessageFromContext(ctx *gin.Context) q.PubSubMessage {
	msg := ctx.MustGet(PubSubMessageCtxKey).(q.PubSubMessage)
	return msg
}
