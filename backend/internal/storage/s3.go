package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"labventroy/backend/internal/apperror"
)

// S3MediaStorage implements MediaStorage using Amazon S3 Pre-signed URLs.
type S3MediaStorage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

// NewS3MediaStorage creates an S3-backed MediaStorage for the target bucket.
func NewS3MediaStorage(client *s3.Client, bucket string) *S3MediaStorage {
	return &S3MediaStorage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        bucket,
	}
}

func (s *S3MediaStorage) GenerateUploadURL(ctx context.Context, key, contentType string, ttl time.Duration) (string, error) {
	req, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", apperror.Internal("failed to generate s3 upload url", err)
	}
	return req.URL, nil
}

func (s *S3MediaStorage) GenerateDownloadURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	req, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", apperror.Internal("failed to generate s3 download url", err)
	}
	return req.URL, nil
}

func (s *S3MediaStorage) DeleteObject(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return apperror.Internal("failed to delete s3 object", err)
	}
	return nil
}
