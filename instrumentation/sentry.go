package instrumentation

import (
	"context"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InitFromEnv reads SENTRY_DSN and SENTRY_ENV from environment, and creates a new client.
func InitFromEnv() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Environment:      os.Getenv("SENTRY_ENV"),
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
	})
}

// NewMiddleware creates new instrumentation miggleware,
// that can be used with gin.
func NewMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}

// Recover tries to recover(), and reports panic.
func Recover(timeout time.Duration) {
	if hub := sentry.CurrentHub(); hub != nil {
		hub.Recover(nil)
		sentry.Flush(timeout)
	}
}

// CopyCtx creates a context.Context containing separate scope,
// which can be safely used in goroutine.
func CopyCtx(ctx context.Context) context.Context {
	if hub := sentry.CurrentHub(); hub != nil {
		clonedHub := hub.Clone()

		switch c := ctx.(type) {
		case *gin.Context:
			return sentry.SetHubOnContext(c.Copy(), clonedHub)
		default:
			return sentry.SetHubOnContext(ctx, clonedHub)
		}
	}
	return ctx
}

// SetTag adds a tag to the current scope.
func SetTag(ctx context.Context, key string, value string) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.Scope().SetTag(key, value)
	}
}

// SetTags assigns multiple tags to the current scope.
func SetTags(ctx context.Context, tags map[string]string) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.Scope().SetTags(tags)
	}
}

// SetUser sets the user for the current scope.
func SetUser(ctx context.Context, id uuid.UUID, email string) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.Scope().SetUser(sentry.User{ID: id.String(), Email: email})
	}
}

// AddBreadcrumb adds new breadcrumb to the current scope.
func AddBreadcrumb(ctx context.Context, category, msg string) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.Scope().AddBreadcrumb(&sentry.Breadcrumb{
			Category: category,
			Message:  msg,
		}, 1000)
	}
}

// StartSpan starts a new span to describe an operation.
func StartSpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	return span
}