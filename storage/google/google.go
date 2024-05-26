package google

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gabriel-vasile/mimetype"
)

type google struct {
	client *storage.Client
	bucket string
}

// NewGoogleStorage creates a new Google Cloud Storage helper
func NewGoogleStorage(client *storage.Client, bucket string) *google {
	return &google{client: client, bucket: bucket}
}

// Upload uploads a file to Google Cloud Storage
func (g *google) Upload(ctx context.Context, name string, r io.ReadSeeker) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	w := g.client.Bucket(g.bucket).Object(name).NewWriter(ctx)

	if contentType, err := mimetype.DetectReader(r); err == nil {
		w.ContentType = contentType.String()
	}

	r.Seek(0, io.SeekStart)

	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}

	return nil
}

// Download downloads a file from Google Cloud Storage. Returns content, content type and error
func (g *google) Download(ctx context.Context, name string) ([]byte, string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	r, err := g.client.Bucket(g.bucket).Object(name).NewReader(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create reader: %w", err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read content: %w", err)
	}

	return b, r.Attrs.ContentType, nil
}

// GetSignedURL returns a signed URL
func (g *google) GetSignedURL(ctx context.Context, name string, expires time.Duration) (string, error) {
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
