package q

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type PubSubQ interface {
	// Enqueue enqueues a task to the Pub/Sub queue.
	Enqueue(ctx context.Context, task *Task, opts ...PubSubOption) (*TaskInfo, error)
}

type pubSubQ struct {
	client *pubsub.Client
}

func NewPubSubQ(client *pubsub.Client) *pubSubQ {
	return &pubSubQ{client: client}
}

func (q *pubSubQ) Enqueue(ctx context.Context, task *Task, opts ...PubSubOption) (*TaskInfo, error) {
	topicID := task.typename

	message := &pubsub.Message{
		Data:       task.payload,
		Attributes: map[string]string{nameKey: task.typename},
	}

	for _, opt := range task.opts {
		switch opt := opt.(type) {
		case uniqueKeyOption:
			message.Attributes[uniqueKeyKey] = string(opt)
		case groupOption:
			topicID = string(opt)
		}
	}

	message.Attributes[groupKey] = topicID

	topic := q.client.Topic(topicID)
	defer topic.Stop()

	for _, opt := range opts {
		switch opt := opt.(type) {
		case orderedKeyOption:
			topic.EnableMessageOrdering = true
			message.OrderingKey = string(opt)
		case orderedByTaskNameOption:
			topic.EnableMessageOrdering = true
			message.OrderingKey = task.typename
		}
	}

	result := topic.Publish(ctx, message)

	messageID, err := result.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue message: %w", err)
	}

	return &TaskInfo{ID: messageID}, nil
}
