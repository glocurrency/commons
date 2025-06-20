package middleware_test

import (
	_ "embed"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	ginfirebasemw "github.com/brokeyourbike/gin-firebase-middleware"
	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/middleware"
	"github.com/glocurrency/commons/router"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/user-email-verified.json
var userEmailVerified []byte

//go:embed testdata/user-email-not-verified.json
var userEmailNotVerified []byte

//go:embed testdata/user-second-factor-phone.json
var userSecondFactorPhone []byte

//go:embed testdata/user-no-email-apple.json
var userNoEmailApple []byte

func TestRequireEmailVerified(t *testing.T) {
	tests := []struct {
		name       string
		header     []byte
		wantStatus int
	}{
		{
			"email not verified",
			userEmailNotVerified,
			http.StatusForbidden,
		},
		{
			"email verified",
			userEmailVerified,
			http.StatusOK,
		},
		{
			"no email, apple.com",
			userNoEmailApple,
			http.StatusOK,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// encoding the header value to match what expected by `ginfirebasemw`
			req.Header.Set("X-Apigateway-Api-Userinfo", base64.RawURLEncoding.EncodeToString(test.header))

			w := httptest.NewRecorder()
			router := router.NewRouterWithValidation()
			router.Use(ginfirebasemw.Middleware())
			router.Use(middleware.RequireEmailVerified())
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "the end.")
			})
			router.ServeHTTP(w, req)

			require.Equal(t, test.wantStatus, w.Code)
		})
	}
}

func TestRequireSecondFactorPhone(t *testing.T) {
	tests := []struct {
		name       string
		header     []byte
		wantStatus int
	}{

		{
			"no second factor",
			userEmailVerified,
			http.StatusForbidden,
		},
		{
			"has second factor phone",
			userSecondFactorPhone,
			http.StatusOK,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// encoding the header value to match what expected by `ginfirebasemw`
			req.Header.Set("X-Apigateway-Api-Userinfo", base64.RawURLEncoding.EncodeToString(test.header))

			w := httptest.NewRecorder()
			router := router.NewRouterWithValidation()
			router.Use(ginfirebasemw.Middleware())
			router.Use(middleware.RequireSecondFactorPhone())
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "the end.")
			})
			router.ServeHTTP(w, req)

			require.Equal(t, test.wantStatus, w.Code)
		})
	}
}

func TestRequireSecondFactor(t *testing.T) {
	tests := []struct {
		name       string
		header     []byte
		wantStatus int
	}{

		{
			"no second factor",
			userEmailVerified,
			http.StatusForbidden,
		},
		{
			"has second factor phone",
			userSecondFactorPhone,
			http.StatusOK,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// encoding the header value to match what expected by `ginfirebasemw`
			req.Header.Set("X-Apigateway-Api-Userinfo", base64.RawURLEncoding.EncodeToString(test.header))

			w := httptest.NewRecorder()
			router := router.NewRouterWithValidation()
			router.Use(ginfirebasemw.Middleware())
			router.Use(middleware.RequireSecondFactor())
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "the end.")
			})
			router.ServeHTTP(w, req)

			require.Equal(t, test.wantStatus, w.Code)
		})
	}
}
