package monitoring_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/glocurrency/commons/monitoring"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
)

func TestNewClient_Err(t *testing.T) {
	os.Unsetenv("NEW_RELIC_APP_NAME")
	os.Unsetenv("NEW_RELIC_LICENSE_KEY")

	client, err := monitoring.NewClient()
	require.Error(t, err)
	require.IsType(t, &newrelic.Application{}, client.Application())
	require.Nil(t, client.Application())
}

func TestNewClient(t *testing.T) {
	os.Setenv("NEW_RELIC_APP_NAME", "dummy")
	os.Setenv("NEW_RELIC_LICENSE_KEY", "1111111111111111111111111111111111111111")

	client, err := monitoring.NewClient()
	require.NoError(t, err)
	require.IsType(t, &newrelic.Application{}, client.Application())
	require.NotNil(t, client.Application())

	ctx := client.StartTransactionContext(context.TODO(), "test")
	defer monitoring.FromContext(ctx).End()
	defer monitoring.StartSegment(ctx, "seg1").End()
	defer monitoring.StartExternalSegment(ctx, &http.Request{}).End()
}
