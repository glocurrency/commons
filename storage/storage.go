package storage

import (
	"context"
	"io"
	"time"
)

type Signer interface {
	// GetSignedURL returns a signed URL
	GetSignedURL(context.Context, string, time.Duration) (string, error)
}

type Uploader interface {
	// Upload uploads a file
	Upload(context.Context, string, io.ReadSeeker) error
}
