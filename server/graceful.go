package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glocurrency/commons/instrumentation"
)

// ServeWithTimeout is a wrapper around http.Server that gracefully shuts down
// the server when a SIGINT/SIGTERM is received OR the provided context is canceled.
func ServeWithTimeout(ctx context.Context, timeout time.Duration, srv *http.Server) {
	go func() {
		instrumentation.NoticeInfo(ctx, "starting server..", instrumentation.WithField("server_addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			instrumentation.NoticeError(ctx, err, "failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// wait for OS signal OR context cancellation
	select {
	case <-quit:
		instrumentation.NoticeInfo(ctx, "shutting down server due to OS signal...")
	case <-ctx.Done():
		instrumentation.NoticeInfo(ctx, "shutting down server due to context cancellation...")
	}

	// derive shutdown context from Background so it isn't affected by parent cancellation
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		instrumentation.NoticeError(ctx, err, "server forced to shutdown")
	}

	instrumentation.NoticeInfo(ctx, "server stopped")
}
