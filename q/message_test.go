package q_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/pubsub.json
var pubsubMsg []byte

//go:embed testdata/pubsub-order.json
var pubsubOrderMsg []byte

//go:embed testdata/pubsub-attributes.json
var pubsubAttributesMsg []byte

func TestPubsubMsg(t *testing.T) {
	t.Parallel()

	var msg q.PubSubMessage
	require.NoError(t, json.Unmarshal(pubsubMsg, &msg))
	assert.Equal(t, "12345", msg.Message.ID)
	assert.Equal(t, "2022-11-15 14:38:33.816 +0000 UTC", msg.Message.PublishTime.String())

	assert.Equal(t, "Runner", string(msg.Message.Data))
	assert.Empty(t, msg.GetUniqueKey())
	assert.Empty(t, msg.GetName())
	assert.Empty(t, msg.GetGroup())
}

func TestPubsubOrderMsg(t *testing.T) {
	t.Parallel()

	var msg q.PubSubMessage
	require.NoError(t, json.Unmarshal(pubsubOrderMsg, &msg))
	assert.Equal(t, "key123", msg.Message.OrderingKey)

	assert.Equal(t, "Runner", string(msg.Message.Data))
	assert.Empty(t, msg.GetUniqueKey())
	assert.Empty(t, msg.GetName())
	assert.Empty(t, msg.GetGroup())
}

func TestMainPubsubAttributesMsg(t *testing.T) {
	t.Parallel()

	var msg q.PubSubMessage
	require.NoError(t, json.Unmarshal(pubsubAttributesMsg, &msg))

	var msgData struct{ Name string }
	assert.NoError(t, msg.UnmarshalData(&msgData))

	assert.Equal(t, "u", msg.GetUniqueKey())
	assert.Equal(t, "n", msg.GetName())
	assert.Equal(t, "g", msg.GetGroup())
}
