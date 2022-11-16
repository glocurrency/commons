package gate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/iamcredentials/v1"
)

type IClient interface {
	// GenerateJWT generates a JWT for the given service account
	GenerateJWT(serviceAccount string, expiry int64) (string, error)
	// AuthenticateRequest authenticates the given request for the given service account
	AuthenticateRequest(req *http.Request, serviceAccount string, expiry int64) error
}

type Client struct {
	iamCredsClient *iamcredentials.Service
}

func NewClient(iamCredsClient *iamcredentials.Service) *Client {
	return &Client{iamCredsClient: iamCredsClient}
}

func (c *Client) GenerateJWT(serviceAccount string, expiry int64) (string, error) {
	now := time.Now().Unix()

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

	return resp.SignedJwt, nil
}

func (c *Client) AuthenticateRequest(req *http.Request, serviceAccount string, expiry int64) error {
	jwt, err := c.GenerateJWT(serviceAccount, expiry)
	if err != nil {
		return fmt.Errorf("error generating jwt: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	return nil
}
