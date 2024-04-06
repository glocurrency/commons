package middleware

import (
	"net/http"

	ginfirebasemw "github.com/brokeyourbike/gin-firebase-middleware"
	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/instrumentation"
)

// RequireIfServiceAccount returns a middleware that checks if the request is coming from a service account
// and if the user ID is the same as the one provided.
func RequireIfServiceAccount(userID string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInfo := ginfirebasemw.GetUserInfo(ctx)

		// skip validation of regular users
		if !userInfo.IsServiceAccount() {
			ctx.Next()
			return
		}

		if userInfo.Sub != userID {
			instrumentation.NoticeWarning(ctx, "service account id do not match",
				instrumentation.WithField("want_id", userID),
				instrumentation.WithField("have_id", userInfo.Sub),
			)
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
