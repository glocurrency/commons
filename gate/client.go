package gate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"google.golang.org/api/iamcredentials/v1"
)

type Client interface {
	// GenerateJWT generates a JWT for the given service account
	GenerateJWT(serviceAccount string, expiry int64) (string, error)
	// AuthenticateRequest authenticates the given request for the given service account
	AuthenticateRequest(req *http.Request, serviceAccount string, expiry int64) error
}

var _ Client = (*client)(nil)

// cachedToken holds the JWT and its absolute expiration time.
type cachedToken struct {
	token     string
	expiresAt int64
}

type client struct {
	iamCredsClient *iamcredentials.Service

	// mu protects the cache map from concurrent read/writes
	mu    sync.RWMutex
	cache map[string]cachedToken
}

func NewClient(iamCredsClient *iamcredentials.Service) *client {
	return &client{
		iamCredsClient: iamCredsClient,
		cache:          make(map[string]cachedToken),
	}
}

func (c *client) GenerateJWT(serviceAccount string, expiry int64) (string, error) {
	now := time.Now().Unix()

	c.mu.RLock()
	cached, exists := c.cache[serviceAccount]
	c.mu.RUnlock()

	// We use a 60-second buffer to ensure the token doesn't expire
	// while the HTTP request is in flight.
	if exists && cached.expiresAt > now+60 {
		return cached.token, nil
	}

	claims := ClaimSet{
		Iat: now,
		Exp: now + expiry,
		Iss: serviceAccount,
		Aud: serviceAccount,
		Sub: serviceAccount,
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("error marshalling claims: %w", err)
	}

	resp, err := c.iamCredsClient.
		Projects.
		ServiceAccounts.
		SignJwt(fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccount), &iamcredentials.SignJwtRequest{Payload: string(payload)}).
		Do()

	if err != nil {
		return "", fmt.Errorf("error signing jwt: %w", err)
	}

	c.mu.Lock()
	c.cache[serviceAccount] = cachedToken{
		token:     resp.SignedJwt,
		expiresAt: now + expiry,
	}
	c.mu.Unlock()

	return resp.SignedJwt, nil
}

func (c *client) AuthenticateRequest(req *http.Request, serviceAccount string, expiry int64) error {
	jwt, err := c.GenerateJWT(serviceAccount, expiry)
	if err != nil {
		return fmt.Errorf("error generating jwt: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	return nil
}
