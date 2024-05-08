package gate_test

import (
	"net/http"
	"testing"

	"github.com/glocurrency/commons/gate"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSimpleClient(t *testing.T) {
	mockClient := gate.NewMockClient(t)
	mockClient.On("GenerateJWT", "service-account", int64(3600)).Return("jwt1", nil)
	mockClient.On("AuthenticateRequest", mock.Anything, "service-account", int64(3600)).Return(nil)

	s := gate.NewSimpleClient(mockClient, "service-account")

	got, err := s.GenerateJWT(3600)
	require.NoError(t, err)
	require.Equal(t, "jwt1", got)

	require.NoError(t, s.AuthenticateRequest(&http.Request{}, 3600))
}
