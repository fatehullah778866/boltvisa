package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSClient struct {
	client     *storage.Client
	bucketName string
}

func NewGCSClient(ctx context.Context, bucketName, credentialsPath string) (*GCSClient, error) {
	var client *storage.Client
	var err error

	if credentialsPath != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	} else {
		// Use default credentials (for Cloud Run/GKE)
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &GCSClient{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (g *GCSClient) UploadFile(ctx context.Context, file io.Reader, objectPath string, contentType string) (string, error) {
	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(objectPath)

	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType
	writer.CacheControl = "public, max-age=3600"

	if _, err := io.Copy(writer, file); err != nil {
		writer.Close()
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Return public URL
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, objectPath)
	return url, nil
}

func (g *GCSClient) DeleteFile(ctx context.Context, objectPath string) error {
	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(objectPath)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func GenerateDocumentPath(applicationID uint, filename string) string {
	timestamp := time.Now().Format("20060102-150405")
	ext := filepath.Ext(filename)
	name := filepath.Base(filename[:len(filename)-len(ext)])
	return fmt.Sprintf("documents/%d/%s-%s%s", applicationID, name, timestamp, ext)
}

func (g *GCSClient) Close() error {
	return g.client.Close()
}
