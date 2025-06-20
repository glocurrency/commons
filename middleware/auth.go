package middleware

import (
	ginfirebasemw "github.com/brokeyourbike/gin-firebase-middleware"
	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/response"
)

func RequireEmailVerified() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInfo := ginfirebasemw.GetUserInfo(ctx)

		// skip validation of service accounts
		if userInfo.IsServiceAccount() {
			ctx.Next()
			return
		}

		if userInfo.Email == "" && userInfo.Firebase.SignInProvider != ginfirebasemw.ProviderPassword {
			ctx.Next()
			return
		}

		if !userInfo.EmailVerified {
			ctx.AbortWithStatusJSON(response.NewErrResponseForbidden("Please verify your email"))
			return
		}

		ctx.Next()
	}
}

func RequireSecondFactorPhone() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInfo := ginfirebasemw.GetUserInfo(ctx)

		// skip validation of service accounts
		if userInfo.IsServiceAccount() {
			ctx.Next()
			return
		}

		if userInfo.Firebase.SignInSecondFactor != ginfirebasemw.SecondFactorPhone {
			ctx.AbortWithStatusJSON(response.NewErrResponseForbidden("Please verify your phone number"))
			return
		}

		ctx.Next()
	}
}

func RequireSecondFactor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInfo := ginfirebasemw.GetUserInfo(ctx)

		// skip validation of service accounts
		if userInfo.IsServiceAccount() {
			ctx.Next()
			return
		}

		if userInfo.Firebase.SignInSecondFactor == "" {
			ctx.AbortWithStatusJSON(response.NewErrResponseForbidden("Please add a second factor authentication"))
			return
		}

		ctx.Next()
	}
}
