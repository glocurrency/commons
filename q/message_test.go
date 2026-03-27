package q_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

// errorReader is a mock io.Reader that forces a read error.
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated body read error")
}

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

	// Test UnmarshalData error handling
	msg.Message.Data = []byte(`{bad-json}`)
	err := msg.UnmarshalData(&msgData)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal payload")
}

// ----------------------------------------------------------------------------
// QMessage Tests
// ----------------------------------------------------------------------------

func TestNewQMessage_CloudTasks(t *testing.T) {
	t.Parallel()

	body := []byte(`{"test": "data"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	// Cloud Tasks specific headers
	req.Header.Set("X-Cloudtasks-Queuename", "my-queue")
	req.Header.Set("nameKey", "task-name")
	req.Header.Set("groupKey", "task-group")
	req.Header.Set("uniqueKey", "task-unique-id")

	msg, err := q.NewQMessage(req)
	require.NoError(t, err)

	assert.Equal(t, "task-name", msg.Name)
	assert.Equal(t, "task-group", msg.Group)
	assert.Equal(t, "task-unique-id", msg.UniqueKey)
	assert.Equal(t, body, msg.Data)

	// Ensure the request body was restored and can be read again if needed
	restoredBody := make([]byte, len(body))
	_, err = req.Body.Read(restoredBody)
	require.NoError(t, err)
	assert.Equal(t, body, restoredBody)
}

func TestNewQMessage_PubSub(t *testing.T) {
	t.Parallel()

	// Reuse the embedded PubSub attributes JSON file
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(pubsubAttributesMsg))

	msg, err := q.NewQMessage(req)
	require.NoError(t, err)

	// In pubsubAttributesMsg, these are "n", "g", and "u" respectively
	assert.Equal(t, "n", msg.Name)
	assert.Equal(t, "g", msg.Group)
	assert.Equal(t, "u", msg.UniqueKey)
	assert.NotEmpty(t, msg.Data)
}

func TestNewQMessage_Errors(t *testing.T) {
	t.Parallel()

	t.Run("Body Read Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", &errorReader{})

		_, err := q.NewQMessage(req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read body")
		assert.Contains(t, err.Error(), "simulated body read error")
	})

	t.Run("PubSub Unmarshal Error", func(t *testing.T) {
		// Valid HTTP request, but invalid JSON for a Pub/Sub event
		badJSON := []byte(`{invalid json payload}`)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(badJSON))

		_, err := q.NewQMessage(req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot unmarshal")
	})
}

func TestQMessage_UnmarshalData(t *testing.T) {
	t.Parallel()

	msg := q.QMessage{
		Data: []byte(`{"id": 42, "status": "active"}`),
	}

	var payload struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}

	t.Run("Success", func(t *testing.T) {
		err := msg.UnmarshalData(&payload)
		require.NoError(t, err)
		assert.Equal(t, 42, payload.ID)
		assert.Equal(t, "active", payload.Status)
	})

	t.Run("Error", func(t *testing.T) {
		msg.Data = []byte(`{bad json}`)
		err := msg.UnmarshalData(&payload)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal payload")
	})
}
