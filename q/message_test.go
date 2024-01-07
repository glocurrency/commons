package q_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/pubsub.json
var pubsubMsg []byte

//go:embed testdata/pubsub-order.json
var pubsubOrderMsg []byte

func TestPubSubMessage(t *testing.T) {
	t.Parallel()

	var msg q.PubSubMessage
	assert.NoError(t, json.Unmarshal(pubsubMsg, &msg))
	assert.Equal(t, "12345", msg.Message.ID)
	assert.Equal(t, "2022-11-15 14:38:33.816 +0000 UTC", msg.Message.PublishTime.String())
}

func TestPubSubMessageOrder(t *testing.T) {
	t.Parallel()

	var msg q.PubSubMessage
	assert.NoError(t, json.Unmarshal(pubsubOrderMsg, &msg))
	assert.Equal(t, "key123", msg.Message.OrderingKey)
}
