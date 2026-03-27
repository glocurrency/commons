package q_test

import (
	"context"
	"net"
	"testing"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// ----------------------------------------------------------------------------
// Mock gRPC Server for Cloud Tasks
// ----------------------------------------------------------------------------

type mockCloudTasksServer struct {
	cloudtaskspb.UnimplementedCloudTasksServer
	req *cloudtaskspb.CreateTaskRequest
	err error
}

func (m *mockCloudTasksServer) CreateTask(ctx context.Context, req *cloudtaskspb.CreateTaskRequest) (*cloudtaskspb.Task, error) {
	m.req = req
	if m.err != nil {
		return nil, m.err
	}

	name := req.Task.Name
	if name == "" {
		name = req.Parent + "/tasks/mock-task-123"
	}

	return &cloudtaskspb.Task{Name: name}, nil
}

func setupMockClient(t *testing.T) (*cloudtasks.Client, *mockCloudTasksServer, func()) {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	srv := &mockCloudTasksServer{}
	cloudtaskspb.RegisterCloudTasksServer(grpcServer, srv)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	client, err := cloudtasks.NewClient(context.Background(),
		option.WithEndpoint(lis.Addr().String()),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	require.NoError(t, err)

	cleanup := func() {
		client.Close()
		grpcServer.Stop()
	}

	return client, srv, cleanup
}

// ----------------------------------------------------------------------------
// Tests
// ----------------------------------------------------------------------------

func TestCloudTasksQ_Enqueue_Marshal(t *testing.T) {
	ps := q.NewCloudTasksQ(q.CloudTasksConfig{}, &cloudtasks.Client{})

	cannotMarshall := make(chan int)

	err := ps.Enqueue(context.TODO(), q.NewTask("test", cannotMarshall))
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to marshal payload")

	err = ps.Enqueue(context.TODO(), nil)
	require.ErrorIs(t, err, q.ErrTaskIsNil)
}

func TestCloudTasksQ_Enqueue_Success(t *testing.T) {
	client, srv, cleanup := setupMockClient(t)
	defer cleanup()

	cfg := q.CloudTasksConfig{
		ProjectID:           "test-project",
		LocationID:          "europe-west1",
		BaseUrl:             "https://example.com/hooks",
		ServiceAccountEmail: "service-account@example.com",
	}

	ps := q.NewCloudTasksQ(cfg, client)
	taskPayload := map[string]string{"msg": "hello"}
	task := q.NewTask("test-task", taskPayload)

	info, err := ps.EnqueueWithInfo(context.Background(), task)
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.Contains(t, info.ID, "mock-task-123")

	require.NotNil(t, srv.req)
	assert.Equal(t, "projects/test-project/locations/europe-west1/queues/test-task", srv.req.Parent)

	httpReq := srv.req.Task.GetHttpRequest()
	require.NotNil(t, httpReq)

	assert.Equal(t, "https://example.com/hooks/test-task", httpReq.Url)
	assert.Equal(t, cloudtaskspb.HttpMethod_POST, httpReq.HttpMethod)
	assert.Equal(t, "application/json", httpReq.Headers["Content-Type"])

	oidcToken := httpReq.GetOidcToken()
	require.NotNil(t, oidcToken)
	assert.Equal(t, "service-account@example.com", oidcToken.ServiceAccountEmail)

	expectedBody := `{"msg":"hello"}`
	assert.JSONEq(t, expectedBody, string(httpReq.Body))
}

func TestCloudTasksQ_Enqueue_WithOptions(t *testing.T) {
	client, srv, cleanup := setupMockClient(t)
	defer cleanup()

	cfg := q.CloudTasksConfig{
		ProjectID:  "test-project",
		LocationID: "europe-west1",
		BaseUrl:    "https://example.com/hooks",
	}

	ps := q.NewCloudTasksQ(cfg, client)

	// Use your new strongly-typed TaskOptions for the task creation
	task := q.NewTask("test-task", nil, q.Group("custom-queue"), q.UniqueKey("dedupe-123"))

	processTime := time.Now().Add(1 * time.Hour)

	// Use your new strongly-typed CloudTasksOption for the Enqueue method
	_, err := ps.EnqueueWithInfo(context.Background(), task, q.ProcessAt(processTime))
	require.NoError(t, err)

	// Verify q.Group changed the queue mapping
	assert.Equal(t, "projects/test-project/locations/europe-west1/queues/custom-queue", srv.req.Parent)

	// Verify q.UniqueKey set a specific Task Name
	assert.Contains(t, srv.req.Task.Name, "projects/test-project/locations/europe-west1/queues/custom-queue/tasks/")

	// Verify q.ProcessAt scheduled the task in the future
	require.NotNil(t, srv.req.Task.ScheduleTime)
	assert.Equal(t, processTime.Unix(), srv.req.Task.ScheduleTime.Seconds)
}

func TestCloudTasksQ_Enqueue_Error(t *testing.T) {
	client, srv, cleanup := setupMockClient(t)
	defer cleanup()

	srv.err = status.Error(codes.Internal, "google cloud is down")

	ps := q.NewCloudTasksQ(q.CloudTasksConfig{}, client)
	task := q.NewTask("test-task", nil)

	_, err := ps.EnqueueWithInfo(context.Background(), task)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to enqueue task")
	assert.Contains(t, err.Error(), "google cloud is down")
}
