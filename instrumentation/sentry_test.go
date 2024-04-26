package instrumentation_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/instrumentation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	defer func() {
		instrumentation.Recover(10 * time.Second)
	}()
}

func TestInitFromEnv(t *testing.T) {
	os.Setenv("SENTRY_DSN", "https://fake@12345.ingest.us.sentry.io/12345")
	os.Setenv("SENTRY_ENV", "dev")

	require.NoError(t, instrumentation.InitFromEnv())
}

func TestInitFromEnvFull(t *testing.T) {
	os.Setenv("SENTRY_DSN", "https://fake@12345.ingest.us.sentry.io/12345")
	os.Setenv("SENTRY_ENV", "dev")
	os.Setenv("SENTRY_TRACE", "true")
	os.Setenv("SENTRY_TRACERATE", "0.1")

	require.NoError(t, instrumentation.InitFromEnv())
}

func TestInitFromEnv_Fail(t *testing.T) {
	os.Setenv("SENTRY_DSN", "not-a-dsn")
	os.Setenv("SENTRY_ENV", "dev")

	require.Error(t, instrumentation.InitFromEnv())
}

func TestMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	app := gin.Default()
	app.Use(instrumentation.NewMiddleware())

	app.GET("/", func(ctx *gin.Context) {
		instrumentation.SetTag(ctx, "a", "b")
		instrumentation.SetTags(ctx, map[string]string{"c": "d"})
		instrumentation.SetUser(ctx, uuid.New(), "john@doe.com")
		instrumentation.AddBreadcrumb(ctx, "auth", "user logged in")

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			instrumentation.AddBreadcrumb(ctx, "routine", "completed")
		}(instrumentation.CopyCtx(ctx))
		wg.Wait()
	})

	app.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
