package q

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrTaskIsNil = fmt.Errorf("task is nil")

type CloudTasksConfig struct {
	ProjectID           string
	LocationID          string
	BaseUrl             string
	ServiceAccountEmail string
}

type CloudTasksQ interface {
	// Enqueue enqueues a task with ClousTasks, and returns info.
	EnqueueWithInfo(ctx context.Context, task *Task, opts ...CloudTasksOption) (*TaskInfo, error)
	// Enqueue enqueues a task with ClousTasks.
	Enqueue(ctx context.Context, task *Task, opts ...CloudTasksOption) error
}

type cloudTasksQ struct {
	cfg    CloudTasksConfig
	client *cloudtasks.Client
}

func NewCloudTasksQ(cfg CloudTasksConfig, client *cloudtasks.Client) *cloudTasksQ {
	return &cloudTasksQ{cfg: cfg, client: client}
}

// Enqueue enqueues a task with ClousTasks.
func (q *cloudTasksQ) Enqueue(ctx context.Context, task *Task, opts ...CloudTasksOption) error {
	_, err := q.EnqueueWithInfo(ctx, task, opts...)
	return err
}

// Enqueue enqueues a task with ClousTasks, and returns info.
func (q *cloudTasksQ) EnqueueWithInfo(ctx context.Context, task *Task, opts ...CloudTasksOption) (info *TaskInfo, err error) {
	if task == nil {
		return nil, ErrTaskIsNil
	}

	queueID := task.typename
	uniqueKey := ""

	var payload []byte

	// marshall payload to JSON
	if task.payload != nil {
		payload, err = json.Marshal(task.payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	for _, opt := range task.opts {
		switch opt := opt.(type) {
		case groupOption:
			queueID = string(opt)
		case uniqueKeyOption:
			uniqueKey = string(opt)
		}
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", q.cfg.ProjectID, q.cfg.LocationID, queueID)

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        fmt.Sprintf("%s/%s", q.cfg.BaseUrl, queueID),
					Body:       payload,
					Headers:    map[string]string{"Content-Type": "application/json", nameKey: task.typename, groupKey: queueID},
					AuthorizationHeader: &cloudtaskspb.HttpRequest_OidcToken{
						OidcToken: &cloudtaskspb.OidcToken{
							ServiceAccountEmail: q.cfg.ServiceAccountEmail,
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
