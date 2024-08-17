package q_test

import (
	"context"
	"testing"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/glocurrency/commons/q"
	"github.com/stretchr/testify/require"
)

func TestCloudTasksQ_Enqueue_Marshal(t *testing.T) {
	cannotMarshall := make(chan int)

	ps := q.NewCloudTasksQ(q.CloudTasksConfig{}, &cloudtasks.Client{})

	err := ps.Enqueue(context.TODO(), q.NewTask("test", cannotMarshall))
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to marshal payload")

	err = ps.Enqueue(context.TODO(), nil)
	require.ErrorIs(t, err, q.ErrTaskIsNil)
}
