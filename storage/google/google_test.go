package google_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/glocurrency/commons/storage/google"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoogleStorage(t *testing.T) {
	bucketName := "test-bucket"

	// 1. Initialize the fake GCS server with some pre-existing data
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

	// 2. Get the real *storage.Client that points to the local fake server
	client := server.Client()
	gcs := google.NewGoogleStorage(client, bucketName)
	ctx := context.Background()

	t.Run("Upload successfully uploads and detects mimetype", func(t *testing.T) {
		fileName := "new-file.html"
		content := []byte("<html><body>Hello World</body></html>")
		reader := bytes.NewReader(content)

		// Execute
		err := gcs.Upload(ctx, fileName, reader)
		require.NoError(t, err)

		// Verify against the fake server directly
		obj, err := server.GetObject(bucketName, fileName)
		require.NoError(t, err)
		assert.Equal(t, content, obj.Content)

		// The mimetype package should have detected this as HTML
		assert.Contains(t, obj.ContentType, "text/html")
	})

	t.Run("Download successfully retrieves file and mimetype", func(t *testing.T) {
		// Execute
		content, contentType, err := gcs.Download(ctx, "existing-file.txt")

		// Verify
		require.NoError(t, err)
		assert.Equal(t, []byte("hello from the fake bucket"), content)
		assert.Equal(t, "text/plain; charset=utf-8", contentType)
	})

	t.Run("Download fails if file does not exist", func(t *testing.T) {
		content, contentType, err := gcs.Download(ctx, "ghost-file.txt")

		require.Error(t, err)
		assert.Nil(t, content)
		assert.Empty(t, contentType)
		assert.Contains(t, err.Error(), "failed to create reader")
	})
}
