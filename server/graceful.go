package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glocurrency/commons/logger"
)

// ServeWithTimeout is a wrapper around http.Server that gracefully shuts down
// the server when a SIGINT or SIGTERM is received.
func ServeWithTimeout(ctx context.Context, timeout time.Duration, srv *http.Server) {
	go func() {
		logger.WithContext(ctx).Debugf("starting HTTP server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithContext(ctx).WithError(err).Fatal("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.WithContext(ctx).Debug("shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithContext(ctx).Debug("server forced to shutdown: ", err)
	}

	logger.WithContext(ctx).Println("server stopped")
}
