package q

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// uniqueKeyKey is the key used to store the unique ID in the message attributes.
const uniqueKeyKey = "uniqueKey"

// nameKey is the key used to store the name in the message attributes.
const nameKey = "nameKey"

// groupKey is the key used to store the group ID in the message attributes.
const groupKey = "groupKey"

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Message struct {
		ID          string            `json:"messageId"`
		PublishTime time.Time         `json:"publishTime"`
		Data        []byte            `json:"data,omitempty"`
		OrderingKey string            `json:"orderingKey,omitempty"`
		Attributes  map[string]string `json:"attributes,omitempty"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func (m *PubSubMessage) GetUniqueKey() string {
	raw, ok := m.Message.Attributes[uniqueKeyKey]
	if !ok {
		return ""
	}
	return raw
}

func (m *PubSubMessage) GetName() string {
	raw, ok := m.Message.Attributes[nameKey]
	if !ok {
		return ""
	}
	return raw
}

func (m *PubSubMessage) GetGroup() string {
	raw, ok := m.Message.Attributes[groupKey]
	if !ok {
		return ""
	}
	return raw
}

func (m *PubSubMessage) UnmarshalData(v interface{}) error {
	if err := json.Unmarshal(m.Message.Data, v); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return nil
}

type QMessage struct {
	Name      string
	Group     string
	UniqueKey string
	Data      []byte
}

func NewQMessage(req *http.Request) (msg QMessage, err error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return msg, fmt.Errorf("cannot read body: %w", err)
	}

	msg.Data = data
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	if req.Header.Get("X-Cloudtasks-Queuename") != "" {
		msg.Name = req.Header.Get(nameKey)
		msg.Group = req.Header.Get(groupKey)
		msg.UniqueKey = req.Header.Get(uniqueKeyKey)
		return msg, nil
	}

	var psm PubSubMessage
	if err := json.Unmarshal(data, &psm); err != nil {
		return msg, fmt.Errorf("cannot unmarshal '%s': %w", string(data), err)
	}

	msg.Name = psm.GetName()
	msg.Group = psm.GetGroup()
	msg.UniqueKey = psm.GetUniqueKey()
	msg.Data = psm.Message.Data
	return msg, nil
}

func (m *QMessage) UnmarshalData(v interface{}) error {
	if err := json.Unmarshal(m.Data, v); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return nil
}
