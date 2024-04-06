package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glocurrency/commons/instrumentation"
)

// ServeWithTimeout is a wrapper around http.Server that gracefully shuts down
// the server when a SIGINT or SIGTERM is received.
func ServeWithTimeout(ctx context.Context, timeout time.Duration, srv *http.Server) {
	go func() {
		instrumentation.NoticeInfo(ctx, "starting server..", instrumentation.WithField("server_addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			instrumentation.NoticeError(ctx, err, "failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	instrumentation.NoticeInfo(ctx, "shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		instrumentation.NoticeError(ctx, err, "server forced to shutdown")
	}

	instrumentation.NoticeInfo(ctx, "server stopped")
}
