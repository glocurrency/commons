package google_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/glocurrency/commons/storage/google" // Adjust to your actual import path
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// errorReadSeeker is a custom mock that fails when Read is called.
// This allows us to simulate an io.Copy failure during the Upload process.
type errorReadSeeker struct{}

func (e *errorReadSeeker) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func (e *errorReadSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func TestGoogleStorage(t *testing.T) {
	bucketName := "test-bucket"

	// 1. Initialize the fake GCS server with pre-existing data
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName:  bucketName,
				Name:        "existing-file.txt",
				ContentType: "text/plain; charset=utf-8",
			},
			Content: []byte("hello from the fake bucket"),
		},
	})
	defer server.Stop()

	// 2. Get the real *storage.Client pointing to the fake server
	client := server.Client()
	gcs := google.NewGoogleStorage(client, bucketName)
	ctx := context.Background()

	// ----------------------------------------------------------------------
	// UPLOAD TESTS
	// ----------------------------------------------------------------------

	t.Run("Upload successfully uploads and detects mimetype", func(t *testing.T) {
		fileName := "new-file.html"
		content := []byte("<html><body>Hello World</body></html>")
		reader := bytes.NewReader(content)

		err := gcs.Upload(ctx, fileName, reader)
		require.NoError(t, err)

		// Verify against the fake server directly
		obj, err := server.GetObject(bucketName, fileName)
		require.NoError(t, err)
		assert.Equal(t, content, obj.Content)
		assert.Contains(t, obj.ContentType, "text/html")
	})

	t.Run("Upload fails on read error (io.Copy failure)", func(t *testing.T) {
		// Pass our custom mock to force a read error
		err := gcs.Upload(ctx, "fail-read.txt", &errorReadSeeker{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to copy content")
		assert.Contains(t, err.Error(), "simulated read error")
	})

	t.Run("Upload fails on upload error (w.Close failure)", func(t *testing.T) {
		// Initialize the helper with a bucket that the fake server doesn't know about
		badGcs := google.NewGoogleStorage(client, "non-existent-bucket")
		reader := bytes.NewReader([]byte("test content"))

		err := badGcs.Upload(ctx, "fail-close.txt", reader)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to upload")
	})

	// ----------------------------------------------------------------------
	// DOWNLOAD TESTS
	// ----------------------------------------------------------------------

	t.Run("Download successfully retrieves file and mimetype", func(t *testing.T) {
		content, contentType, err := gcs.Download(ctx, "existing-file.txt")

		require.NoError(t, err)
		assert.Equal(t, []byte("hello from the fake bucket"), content)
		assert.Equal(t, "text/plain; charset=utf-8", contentType)
	})

	t.Run("Download fails if file does not exist (NewReader failure)", func(t *testing.T) {
		content, contentType, err := gcs.Download(ctx, "ghost-file.txt")

		require.Error(t, err)
		assert.Nil(t, content)
		assert.Empty(t, contentType)
		assert.Contains(t, err.Error(), "failed to create reader")
	})

	// Note: Triggering an `io.ReadAll` error on a successful `NewReader` is highly
	// circumstantial (e.g., connection drop mid-stream). The NewReader failure
	// above sufficiently covers the primary failure state for this method.

	// ----------------------------------------------------------------------
	// SIGNED URL TESTS
	// ----------------------------------------------------------------------

	t.Run("GetSignedURL fails without configured credentials", func(t *testing.T) {
		// The fake-gcs-server client is initialized without authentication.
		// Signed URLs require a private key loaded into the client to succeed.
		// We use this constraint to successfully test your error wrapping logic.
		url, err := gcs.GetSignedURL(ctx, "existing-file.txt", time.Hour)

		require.Error(t, err)
		assert.Empty(t, url)
		assert.Contains(t, err.Error(), "failed to get signed URL")
	})
}
