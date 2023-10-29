package google

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gabriel-vasile/mimetype"
	"github.com/glocurrency/commons/monitoring"
)

type uploader struct {
	client *storage.Client
	bucket string
}

// NewUploader creates a new Google uploader
func NewUploader(client *storage.Client, bucket string) *uploader {
	return &uploader{client: client, bucket: bucket}
}

// Upload uploads a file to Google Cloud Storage
func (g *uploader) Upload(ctx context.Context, name string, r io.ReadSeeker) error {
	defer monitoring.StartSegment(ctx, "google:uploader:Upload").End()

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
