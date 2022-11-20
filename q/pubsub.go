package q

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type IPubSubQ interface {
	// Enqueue enqueues a task to the Pub/Sub queue.
	Enqueue(ctx context.Context, task *Task, opts ...PubSubOption) (*TaskInfo, error)
}

type PubSubQ struct {
	client *pubsub.Client
}

func NewPubSubQ(client *pubsub.Client) *PubSubQ {
	return &PubSubQ{client: client}
}

func (q *PubSubQ) Enqueue(ctx context.Context, task *Task, opts ...PubSubOption) (*TaskInfo, error) {
	topicID := task.typename

	message := &pubsub.Message{
		Data:       task.payload,
		Attributes: map[string]string{},
	}

	for _, opt := range task.opts {
		switch opt := opt.(type) {
		case uniqueKeyOption:
			message.Attributes[uniqueKeyKey] = string(opt)
		}
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case orderedKeyOption:
			message.OrderingKey = string(opt)
		case orderedByTaskNameOption:
			message.OrderingKey = task.typename
		case topicOption:
			topicID = string(opt)
		}
	}

	topic := q.client.Topic(topicID)
	defer topic.Stop()

	// no harm in always enabling this feature
	// since it depends if OrderingKey is empty or not.
	topic.EnableMessageOrdering = true

	result := topic.Publish(ctx, message)

	messageID, err := result.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue message: %w", err)
	}

	return &TaskInfo{ID: messageID}, nil
}
