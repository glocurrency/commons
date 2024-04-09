package q_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/require"
)

func TestEnqueue_Marshal(t *testing.T) {
	cannotMarshall := make(chan int)

	ps := q.NewPubSubQ(&pubsub.Client{})

	err := ps.Enqueue(context.TODO(), q.NewTask("test", cannotMarshall))
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to marshal payload")
}
