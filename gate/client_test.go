package gate_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/glocurrency/commons/gate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
)

// mockTransport implements http.RoundTripper to intercept HTTP calls
type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

// setupMockService is a helper to create an iamcredentials.Service with a mocked HTTP client
func setupMockService(t *testing.T, rtFunc func(req *http.Request) (*http.Response, error)) *iamcredentials.Service {
	httpClient := &http.Client{
		Transport: &mockTransport{roundTripFunc: rtFunc},
	}

	svc, err := iamcredentials.NewService(context.Background(), option.WithHTTPClient(httpClient))
	require.NoError(t, err)
	return svc
}

func TestNewClient(t *testing.T) {
	c := gate.NewClient(&iamcredentials.Service{})
	require.NotNil(t, c)
}

func TestClient_GenerateJWT(t *testing.T) {
	serviceAccount := "test-sa@project.iam.gserviceaccount.com"
	expectedJWT := "header.payload.signature"

	t.Run("success", func(t *testing.T) {
		mockSvc := setupMockService(t, func(req *http.Request) (*http.Response, error) {
			// Mock the exact JSON response the Google API would return
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
		jwt, err := client.GenerateJWT(serviceAccount, 3600)

		assert.NoError(t, err)
		assert.Equal(t, expectedJWT, jwt)
	})

	t.Run("api error", func(t *testing.T) {
		mockSvc := setupMockService(t, func(req *http.Request) (*http.Response, error) {
			// Simulate a network or API failure
			return nil, errors.New("simulated network timeout")
		})

		client := gate.NewClient(mockSvc)
		jwt, err := client.GenerateJWT(serviceAccount, 3600)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error signing jwt")
		assert.Empty(t, jwt)
	})
}

func TestClient_AuthenticateRequest(t *testing.T) {
	serviceAccount := "test-sa@project.iam.gserviceaccount.com"
	expectedJWT := "header.payload.signature"

	t.Run("success", func(t *testing.T) {
		mockSvc := setupMockService(t, func(req *http.Request) (*http.Response, error) {
			respBody := &iamcredentials.SignJwtResponse{SignedJwt: expectedJWT}
			b, _ := json.Marshal(respBody)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(b)),
				Header:     make(http.Header),
			}, nil
		})

		client := gate.NewClient(mockSvc)
		req, err := http.NewRequest(http.MethodGet, "https://api.example.com/data", nil)
		require.NoError(t, err)

		err = client.AuthenticateRequest(req, serviceAccount, 3600)

		assert.NoError(t, err)
		assert.Equal(t, "Bearer "+expectedJWT, req.Header.Get("Authorization"))
	})

	t.Run("fails when jwt generation fails", func(t *testing.T) {
		mockSvc := setupMockService(t, func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("simulated network timeout")
		})

		client := gate.NewClient(mockSvc)
		req, err := http.NewRequest(http.MethodGet, "https://api.example.com/data", nil)
		require.NoError(t, err)

		err = client.AuthenticateRequest(req, serviceAccount, 3600)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error generating jwt")
		assert.Empty(t, req.Header.Get("Authorization"))
	})
}
