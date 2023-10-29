package middleware

import (
	ginfirebasemw "github.com/brokeyourbike/gin-firebase-middleware"
	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/logger"
	"github.com/glocurrency/commons/response"
	"github.com/google/uuid"
)

const UserUUIDCtxKey = "adminCtx"

func UserUUIDCtx() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raw := ginfirebasemw.GetUserID(ctx)

		id, err := uuid.Parse(raw)
		if err != nil {
			logger.WithContext(ctx).
				WithError(err).
				WithField("user_id", raw).
				Error("admin id cannot be parsed")

			ctx.AbortWithStatusJSON(response.NewErrResponseException("User ID not valid UUID"))
			return
		}

		ctx.Set(UserUUIDCtxKey, id)
		ctx.Next()
	}
}

// MustGetUserUUIDFromContext returns the user ID from the context.
func MustGetUserUUIDFromContext(ctx *gin.Context) uuid.UUID {
	return ctx.MustGet(UserUUIDCtxKey).(uuid.UUID)
}
