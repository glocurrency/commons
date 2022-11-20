package q

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message interface {
	GetUniqueKey() string
	UnmarshalData(v interface{}) error
}

// uniqueKeyKey is the key used to store the unique ID in the message attributes.
const uniqueKeyKey = "uniqueKey"

// topicKey is the key used to store the topic ID in the message attributes.
const topicKey = "topicKey"

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

func (m *PubSubMessage) GetTopicKey() string {
	raw, ok := m.Message.Attributes[topicKey]
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
