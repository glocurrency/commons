package gate_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/glocurrency/commons/gate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
)

// mockTransport implements http.RoundTripper and counts API calls.
// We added a mutex here to ensure thread-safe counting if you run tests in parallel.
type mockTransport struct {
	mu            sync.Mutex
	callCount     int
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	m.callCount++
	m.mu.Unlock()
	return m.roundTripFunc(req)
}

func (m *mockTransport) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

// setupMockService now returns both the service and the transport so we can inspect the call count.
func setupMockService(t *testing.T, rtFunc func(req *http.Request) (*http.Response, error)) (*iamcredentials.Service, *mockTransport) {
	transport := &mockTransport{roundTripFunc: rtFunc}
	httpClient := &http.Client{
		Transport: transport,
	}

	svc, err := iamcredentials.NewService(context.Background(), option.WithHTTPClient(httpClient))
	require.NoError(t, err)
	return svc, transport
}

func TestClient_CacheBehavior(t *testing.T) {
	serviceAccount1 := "sa-one@project.iam.gserviceaccount.com"
	serviceAccount2 := "sa-two@project.iam.gserviceaccount.com"
	expectedJWT := "header.payload.signature"

	mockSvc, transport := setupMockService(t, func(req *http.Request) (*http.Response, error) {
		respBody := &iamcredentials.SignJwtResponse{
			KeyId:     "some-key-id",
			SignedJwt: expectedJWT,
		}
		b, _ := json.Marshal(respBody)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(b)),
			Header:     make(http.Header),
		}, nil
	})

	client := gate.NewClient(mockSvc)

	t.Run("first call makes an API request", func(t *testing.T) {
		jwt, err := client.GenerateJWT(serviceAccount1, 3600)
		assert.NoError(t, err)
		assert.Equal(t, expectedJWT, jwt)
		assert.Equal(t, 1, transport.GetCallCount())
	})

	t.Run("subsequent calls for the same SA use the cache", func(t *testing.T) {
		// Call it 5 more times
		for i := 0; i < 5; i++ {
			jwt, err := client.GenerateJWT(serviceAccount1, 3600)
			assert.NoError(t, err)
			assert.Equal(t, expectedJWT, jwt)
		}
		// The call count should STILL be 1!
		assert.Equal(t, 1, transport.GetCallCount())
	})

	t.Run("call for a different SA triggers a new API request", func(t *testing.T) {
		jwt, err := client.GenerateJWT(serviceAccount2, 3600)
		assert.NoError(t, err)
		assert.Equal(t, expectedJWT, jwt)

		// The call count should now increment to 2
		assert.Equal(t, 2, transport.GetCallCount())
	})
}

func TestClient_AuthenticateRequest_Fails(t *testing.T) {
	mockSvc, _ := setupMockService(t, func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("simulated network timeout")
	})

	client := gate.NewClient(mockSvc)
	req, err := http.NewRequest(http.MethodGet, "https://api.example.com/data", nil)
	require.NoError(t, err)

	err = client.AuthenticateRequest(req, "bad-sa@project", 3600)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error generating jwt")
}
