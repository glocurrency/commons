package server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/glocurrency/commons/server"
)

func TestServeWithTimeout_GracefulShutdownViaContext(t *testing.T) {
	// 1. Setup a dummy HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Note: using port ":0" asks the OS to assign a random free port,
	// preventing port collisions if tests run in parallel.
	srv := &http.Server{
		Addr:    "127.0.0.1:0",
		Handler: handler,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // fallback cleanup

	// 2. Run the server in a goroutine
	done := make(chan struct{})
	go func() {
		server.ServeWithTimeout(ctx, 2*time.Second, srv)
		close(done)
	}()

	// Give the server a tiny fraction of a second to bind to the port
	time.Sleep(50 * time.Millisecond)

	// 3. Trigger shutdown
	cancel()

	// 4. Verify the server actually shuts down and exits the function
	select {
	case <-done:
		// Success! The function returned.
	case <-time.After(3 * time.Second):
		// If it takes longer than our 2s timeout + 1s buffer, something is deadlocked.
		t.Fatal("ServeWithTimeout did not shut down in time")
	}
}
