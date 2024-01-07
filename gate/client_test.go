package gate_test

import (
	"testing"

	"github.com/glocurrency/commons/gate"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/iamcredentials/v1"
)

func TestNewClient(t *testing.T) {
	c := gate.NewClient(&iamcredentials.Service{})
	require.NotNil(t, c)
}
