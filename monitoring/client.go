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

func (c Client) Application() *newrelic.Application {
	return c.app
}

func (c Client) StartTransaction(name string) *newrelic.Transaction {
	return c.app.StartTransaction(name)
}

func (c Client) StartTransactionContext(ctx context.Context, name string) (context.Context, *newrelic.Transaction) {
	trx := c.StartTransaction(name)
	return newrelic.NewContext(ctx, trx), trx
}

func StartSegment(ctx context.Context, name string) *newrelic.Segment {
	return newrelic.FromContext(ctx).StartSegment(name)
}

func StartExternalSegment(ctx context.Context, req *http.Request) *newrelic.ExternalSegment {
	return newrelic.StartExternalSegment(newrelic.FromContext(ctx), req)
}

func FromContext(ctx context.Context) *newrelic.Transaction {
	return newrelic.FromContext(ctx)
}
