// Package objectstorage provides a thin S3-compatible client for services
// like RustFS, MinIO, AWS S3, Cloudflare R2, and Backblaze B2.
//
// Built on aws-sdk-go-v2 (Apache 2.0).
package objectstorage

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config holds the connection parameters for an S3-compatible object store.
type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Region    string
}

// Client wraps an S3-compatible client.
type Client struct {
	s3 *s3.Client
}

// New creates a new Client connected to an S3-compatible endpoint.
// Passing a context is required to honor cancellation during credential
// loading; use context.Background() when no deadline is needed.
func New(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("objectstorage: load aws config: %w", err)
	}

	endpoint := cfg.Endpoint
	if cfg.UseSSL && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}
	if !cfg.UseSSL && !strings.HasPrefix(endpoint, "http://") {
		endpoint = "http://" + endpoint
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &Client{s3: client}, nil
}

// RemoveObject deletes an object from a bucket.
func (c *Client) RemoveObject(ctx context.Context, bucket, key string) error {
	_, err := c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("objectstorage: remove %s/%s: %w", bucket, key, err)
	}
	return nil
}

// PresignedPutObject returns a pre-signed URL suitable for uploading an
// object via HTTP PUT. The caller is responsible for setting Content-Type
// and Content-Length headers when using the returned URL.
func (c *Client) PresignedPutObject(ctx context.Context, bucket, key string, expires time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.s3)

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", fmt.Errorf("objectstorage: presign put %s/%s: %w", bucket, key, err)
	}
	return req.URL, nil
}

// ParseURL extracts bucket and key from a full object URL.
//
// Examples:
//
//	https://s3.example.com/bucket/key/nested/file.png  → ("bucket", "key/nested/file.png")
//	http://localhost:9000/bucket/file.png               → ("bucket", "file.png")
func ParseURL(rawURL string) (bucket, key string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("objectstorage: parse url: %w", err)
	}
	parts := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("objectstorage: path too short, expected /bucket/key: %s", u.Path)
	}
	return parts[0], parts[1], nil
}
