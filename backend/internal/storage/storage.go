package storage

import (
	"context"
	"time"
)

// MediaStorage defines object storage operations for component photos.
type MediaStorage interface {
	GenerateUploadURL(ctx context.Context, key, contentType string, ttl time.Duration) (string, error)
	GenerateDownloadURL(ctx context.Context, key string, ttl time.Duration) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

// NoopMediaStorage is a no-op MediaStorage used when S3 is not configured (e.g., local demo or unit tests).
type NoopMediaStorage struct{}

func (NoopMediaStorage) GenerateUploadURL(_ context.Context, key, _ string, _ time.Duration) (string, error) {
	return "https://example.local/upload/" + key, nil
}

func (NoopMediaStorage) GenerateDownloadURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return "https://example.local/download/" + key, nil
}

func (NoopMediaStorage) DeleteObject(_ context.Context, _ string) error {
	return nil
}
