package google

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/glocurrency/commons/monitoring"
)

type signer struct {
	client *storage.Client
	bucket string
}

// NewSigner creates a new Google signer
func NewSigner(client *storage.Client, bucket string) *signer {
	return &signer{client: client, bucket: bucket}
}

// GetSignedURL returns a signed URL
func (g *signer) GetSignedURL(ctx context.Context, name string, expires time.Duration) (string, error) {
	defer monitoring.StartSegment(ctx, "google:signer:GetSignedURL").End()

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expires),
	}

	u, err := g.client.Bucket(g.bucket).SignedURL(name, opts)
	if err != nil {
		return "", fmt.Errorf("failed to get signed URL: %w", err)
	}

	return u, nil
}
