package gate

import (
	"net/http"
)

type SimpleClient interface {
	GenerateJWT(expiry int64) (string, error)
	AuthenticateRequest(req *http.Request, expiry int64) error
}

var _ SimpleClient = (*simpleClient)(nil)

type simpleClient struct {
	client         Client
	serviceAccount string
}

func NewSimpleClient(client Client, serviceAccount string) *simpleClient {
	return &simpleClient{client: client, serviceAccount: serviceAccount}
}

func (c *simpleClient) GenerateJWT(expiry int64) (string, error) {
	return c.client.GenerateJWT(c.serviceAccount, expiry)
}

func (c *simpleClient) AuthenticateRequest(req *http.Request, expiry int64) error {
	return c.client.AuthenticateRequest(req, c.serviceAccount, expiry)
}
