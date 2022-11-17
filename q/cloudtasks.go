package q

import (
	"context"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

type Config interface {
	GetProjectID() string
	GetLocationID() string
	GetBaseUrl() string
}

type CloudTasksQ struct {
	cfg    Config
	client *cloudtasks.Client
}

func NewCloudTasksQ(cfg Config, client *cloudtasks.Client) *CloudTasksQ {
	return &CloudTasksQ{cfg: cfg, client: client}
}

func (q *CloudTasksQ) Enqueue(ctx context.Context, task *Task, opts ...CloudTasksOption) (*TaskInfo, error) {
	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", q.cfg.GetProjectID(), q.cfg.GetLocationID(), task.typename)

	// Build the Task payload.
	req := &cloudtaskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        fmt.Sprintf("%s/%s", q.cfg.GetBaseUrl(), task.typename),
					Body:       task.payload,
				},
			},
		},
	}

	// TODO: use headers to send attributes
	// TODO: use common interface for the messages received from pubsub and cloudtasks, to contain UniqueKey

	t, err := q.client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue task: %w", err)
	}

	return &TaskInfo{ID: t.GetName()}, nil
}
