package monitoring_test

import (
	"context"
	"os"
	"testing"

	"github.com/glocurrency/commons/monitoring"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	os.Setenv("NEW_RELIC_APP_NAME", "dummy")
	os.Setenv("NEW_RELIC_LICENSE_KEY", "1111111111111111111111111111111111111111")

	client, err := monitoring.NewClient()
	assert.NoError(t, err)
	assert.IsType(t, &newrelic.Application{}, client.Application())

	ctx := client.StartTransactionContext(context.TODO(), "test")
	defer monitoring.FromContext(ctx).End()
}
