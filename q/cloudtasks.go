package q

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Config interface {
	GetProjectID() string
	GetLocationID() string
	GetBaseUrl() string
	GetServiceAccountEmail() string
}

type ICloudTasksQ interface {
	// Enqueue enqueues a task to the ClousTasks queue.
	Enqueue(ctx context.Context, task *Task, opts ...CloudTasksOption) (*TaskInfo, error)
}

type CloudTasksQ struct {
	cfg    Config
	client *cloudtasks.Client
}

func NewCloudTasksQ(cfg Config, client *cloudtasks.Client) *CloudTasksQ {
	return &CloudTasksQ{cfg: cfg, client: client}
}

func (q *CloudTasksQ) Enqueue(ctx context.Context, task *Task, opts ...CloudTasksOption) (*TaskInfo, error) {
	queueID := task.typename
	uniqueKey := ""

	for _, opt := range task.opts {
		switch opt := opt.(type) {
		case groupOption:
			queueID = string(opt)
		case uniqueKeyOption:
			uniqueKey = string(opt)
		}
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", q.cfg.GetProjectID(), q.cfg.GetLocationID(), queueID)

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        fmt.Sprintf("%s/%s", q.cfg.GetBaseUrl(), queueID),
					Body:       task.payload,
					Headers:    map[string]string{"Content-Type": "application/json", nameKey: task.typename, groupKey: queueID},
					AuthorizationHeader: &cloudtaskspb.HttpRequest_OidcToken{
						OidcToken: &cloudtaskspb.OidcToken{
							ServiceAccountEmail: q.cfg.GetServiceAccountEmail(),
						},
					},
				},
			},
		},
	}

	if uniqueKey != "" {
		req.Task.GetHttpRequest().Headers[uniqueKeyKey] = uniqueKey
		req.Task.Name = fmt.Sprintf("%s/tasks/%x", queuePath, sha256.Sum256([]byte(uniqueKey)))
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case processAtOption:
			req.Task.ScheduleTime = timestamppb.New(time.Time(opt))
		}
	}

	t, err := q.client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue task: %w", err)
	}

	return &TaskInfo{ID: t.GetName()}, nil
}
