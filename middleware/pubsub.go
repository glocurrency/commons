package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/logger"
	"github.com/glocurrency/commons/q"
)

const PubSubMessageCtx = "pubSubMessageCtx"

type Locker interface {
	TryToLock(ctx context.Context, key string) error
}

type pubSubCtx struct {
	locker Locker
}

func NewPubSubCtx(l Locker) *pubSubCtx {
	return &pubSubCtx{locker: l}
}

func (m *pubSubCtx) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var msg q.PubSubMessage
		if err := ctx.ShouldBindJSON(&msg); err != nil {
			logger.WithContext(ctx).WithError(err).Error("cannot unmarshal pubsub message")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		logger.WithContext(ctx).WithField("pubsub_message", msg).Debug("pubsub message received")

		if msg.GetUniqueKey() != "" {
			if err := m.locker.TryToLock(ctx, msg.GetUniqueKey()); err != nil {
				logger.WithContext(ctx).
					WithError(err).
					WithField("unique_key", msg.GetUniqueKey()).
					Error("cannot lock task")

				ctx.AbortWithStatus(http.StatusUnprocessableEntity)
				return
			}
		}

		ctx.Set(PubSubMessageCtx, msg)
		ctx.Next()
	}
}

// MustGetMessageFromContext returns the PubSub message from the context.
func MustGetMessageFromContext(ctx *gin.Context) q.PubSubMessage {
	msg := ctx.MustGet(PubSubMessageCtx).(q.PubSubMessage)
	return msg
}
