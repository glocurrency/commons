package monitoring

import (
	"context"
	"fmt"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type Client struct {
	app *newrelic.Application
}

func NewClient() (Client, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigAppLogDecoratingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(cfg *newrelic.Config) {
			cfg.ErrorCollector.RecordPanics = true
		},
	)

	if err != nil {
		return Client{}, fmt.Errorf("cannot create newrelic app: %w", err)
	}

	return Client{app: app}, nil
}

// Application returns the newrelic application.
func (c Client) Application() *newrelic.Application {
	return c.app
}

// StartTransaction starts a newrelic transaction.
func (c Client) StartTransaction(name string) *newrelic.Transaction {
	return c.app.StartTransaction(name)
}

// StartTransactionContext starts a newrelic transaction and puts it in the context.
func (c Client) StartTransactionContext(ctx context.Context, name string) context.Context {
	return newrelic.NewContext(ctx, c.StartTransaction(name))
}

// StartSegment starts a newrelic segment. It extracts the transaction from the context.
func StartSegment(ctx context.Context, name string) *newrelic.Segment {
	return newrelic.FromContext(ctx).StartSegment(name)
}

// StartExternalSegment starts a newrelic external segment. It extracts the transaction from the context.
func StartExternalSegment(ctx context.Context, req *http.Request) *newrelic.ExternalSegment {
	return newrelic.StartExternalSegment(newrelic.FromContext(ctx), req)
}

// FromContext extracts the transaction from the context.
func FromContext(ctx context.Context) *newrelic.Transaction {
	return newrelic.FromContext(ctx)
}
