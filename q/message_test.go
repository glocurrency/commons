package q_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPubSubMessage(t *testing.T) {
	t.Parallel()

	payload, err := os.ReadFile("testdata/pubsub.json")
	require.NoError(t, err)

	var msg q.PubSubMessage
	err = json.Unmarshal(payload, &msg)
	assert.NoError(t, err)
	assert.Equal(t, "12345", msg.Message.ID)
	assert.Equal(t, "2022-11-15 14:38:33.816 +0000 UTC", msg.Message.PublishTime.String())
}

func TestPubSubMessageOrder(t *testing.T) {
	t.Parallel()

	payload, err := os.ReadFile("testdata/pubsub-order.json")
	require.NoError(t, err)

	var msg q.PubSubMessage
	err = json.Unmarshal(payload, &msg)
	assert.NoError(t, err)
	assert.Equal(t, "key123", msg.Message.OrderingKey)
}
