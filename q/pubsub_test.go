package q_test

import (
	"context"
	"net"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
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
// Mock gRPC Server for Pub/Sub
// ----------------------------------------------------------------------------

// mockPublisherServer implements the gRPC interface for Pub/Sub Publishers.
type mockPublisherServer struct {
	pubsubpb.UnimplementedPublisherServer
	req *pubsubpb.PublishRequest
	err error
}

func (m *mockPublisherServer) Publish(ctx context.Context, req *pubsubpb.PublishRequest) (*pubsubpb.PublishResponse, error) {
	m.req = req
	if m.err != nil {
		return nil, m.err
	}
	// Return a fake generated message ID on success
	return &pubsubpb.PublishResponse{MessageIds: []string{"mock-msg-123"}}, nil
}

// setupMockPubSubClient spins up an in-memory gRPC server and returns a Pub/Sub client connected to it.
func setupMockPubSubClient(t *testing.T) (*pubsub.Client, *mockPublisherServer, func()) {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	srv := &mockPublisherServer{}
	pubsubpb.RegisterPublisherServer(grpcServer, srv)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	client, err := pubsub.NewClient(context.Background(), "test-project",
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

func TestPubSubQ_Enqueue_Marshal(t *testing.T) {
	// Doesn't need the mock server since it fails before network calls
	ps := q.NewPubSubQ(&pubsub.Client{})

	cannotMarshall := make(chan int)

	err := ps.Enqueue(context.TODO(), q.NewTask("test", cannotMarshall))
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to marshal payload")

	err = ps.Enqueue(context.TODO(), nil)
	require.ErrorIs(t, err, q.ErrTaskIsNil)
}

func TestPubSubQ_Enqueue_Success(t *testing.T) {
	client, srv, cleanup := setupMockPubSubClient(t)
	defer cleanup()

	ps := q.NewPubSubQ(client)
	taskPayload := map[string]string{"msg": "hello"}
	task := q.NewTask("test-task", taskPayload)

	info, err := ps.EnqueueWithInfo(context.Background(), task)

	require.NoError(t, err)
	require.NotNil(t, info)
	assert.Equal(t, "mock-msg-123", info.ID)

	// Verify the outgoing request to the mock server
	require.NotNil(t, srv.req)

	// Check the target topic
	assert.Equal(t, "projects/test-project/topics/test-task", srv.req.Topic)
	require.Len(t, srv.req.Messages, 1)

	msg := srv.req.Messages[0]

	// Check Payload
	expectedBody := `{"msg":"hello"}`
	assert.JSONEq(t, expectedBody, string(msg.Data))

	// Check Base Attributes
	assert.Equal(t, "test-task", msg.Attributes["nameKey"])
	assert.Equal(t, "test-task", msg.Attributes["groupKey"]) // Defaults to task.typename
}

func TestPubSubQ_Enqueue_TaskOptions(t *testing.T) {
	client, srv, cleanup := setupMockPubSubClient(t)
	defer cleanup()

	ps := q.NewPubSubQ(client)

	// Applying task options (Group & UniqueKey)
	task := q.NewTask("test-task", nil, q.Group("custom-group"), q.UniqueKey("dedupe-abc"))

	_, err := ps.EnqueueWithInfo(context.Background(), task)
	require.NoError(t, err)
	require.NotNil(t, srv.req)
	require.Len(t, srv.req.Messages, 1)

	// Group Option should override the target Topic and the groupKey attribute
	assert.Equal(t, "projects/test-project/topics/custom-group", srv.req.Topic)
	assert.Equal(t, "custom-group", srv.req.Messages[0].Attributes["groupKey"])

	// UniqueKey Option should set the uniqueKey attribute
	assert.Equal(t, "dedupe-abc", srv.req.Messages[0].Attributes["uniqueKey"])
}

func TestPubSubQ_Enqueue_PubSubOptions(t *testing.T) {
	client, srv, cleanup := setupMockPubSubClient(t)
	defer cleanup()

	ps := q.NewPubSubQ(client)
	task := q.NewTask("test-task", nil)

	t.Run("OrderedKey Option", func(t *testing.T) {
		_, err := ps.EnqueueWithInfo(context.Background(), task, q.OrderedKey("order-key-1"))
		require.NoError(t, err)
		require.NotNil(t, srv.req)

		msg := srv.req.Messages[0]
		assert.Equal(t, "order-key-1", msg.OrderingKey)
	})

	t.Run("OrderedByTaskName Option", func(t *testing.T) {
		_, err := ps.EnqueueWithInfo(context.Background(), task, q.OrderedByTaskName())
		require.NoError(t, err)
		require.NotNil(t, srv.req)

		msg := srv.req.Messages[0]
		// Should default the OrderingKey to the task name
		assert.Equal(t, "test-task", msg.OrderingKey)
	})
}

func TestPubSubQ_Enqueue_Error(t *testing.T) {
	client, srv, cleanup := setupMockPubSubClient(t)
	defer cleanup()

	srv.err = status.Error(codes.NotFound, "pubsub topic not found")

	ps := q.NewPubSubQ(client)
	task := q.NewTask("test-task", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := ps.EnqueueWithInfo(ctx, task)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to enqueue message")
	assert.Contains(t, err.Error(), "pubsub topic not found") // Update this to match
}
