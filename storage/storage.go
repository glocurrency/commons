package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Signer
	Uploader
	Downloaded
}

type Signer interface {
	// GetSignedURL returns a signed URL. The URL is valid for the given duration
	GetSignedURL(context.Context, string, time.Duration) (string, error)
}

type Uploader interface {
	// Upload uploads a file. Returns an error if the upload fails
	Upload(context.Context, string, io.ReadSeeker) error
}

type Downloaded interface {
	// Download downloads a file. Returns content, content type and error
	Download(ctx context.Context, name string) ([]byte, string, error)
}
