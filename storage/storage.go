package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Signer
	Uploader
	Downloader
}

type Signer interface {
	// GetSignedURL returns a signed URL, and an error if the operation fails.
	GetSignedURL(context.Context, string, time.Duration) (string, error)
}

type Uploader interface {
	// Upload uploads a file.
	// Returns an error if the upload fails.
	Upload(context.Context, string, io.ReadSeeker) error
}

type Downloader interface {
	// Download downloads a file.
	// Returns content, content type and error if the download fails.
	Download(ctx context.Context, name string) ([]byte, string, error)
}
