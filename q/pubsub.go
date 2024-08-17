package q

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type PubSubQ interface {
	// Enqueue enqueues a task with Pub/Sub, and returns info.
	EnqueueWithInfo(ctx context.Context, task *Task, opts ...PubSubOption) (*TaskInfo, error)
	// Enqueue enqueues a task with Pub/Sub.
	Enqueue(ctx context.Context, task *Task, opts ...PubSubOption) error
}

type pubSubQ struct {
	client *pubsub.Client
}

func NewPubSubQ(client *pubsub.Client) *pubSubQ {
	return &pubSubQ{client: client}
}

// Enqueue enqueues a task with Pub/Sub.
func (q *pubSubQ) Enqueue(ctx context.Context, task *Task, opts ...PubSubOption) error {
	_, err := q.EnqueueWithInfo(ctx, task, opts...)
	return err
}

// Enqueue enqueues a task with Pub/Sub, and returns info.
func (q *pubSubQ) EnqueueWithInfo(ctx context.Context, task *Task, opts ...PubSubOption) (info *TaskInfo, err error) {
	if task == nil {
		return nil, ErrTaskIsNil
	}

	topicID := task.typename

	var payload []byte

	// marshall payload to JSON
	if task.payload != nil {
		payload, err = json.Marshal(task.payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	message := &pubsub.Message{
		Data:       payload,
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
