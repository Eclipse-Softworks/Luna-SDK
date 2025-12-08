package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// BucketsResource provides access to bucket operations
type BucketsResource struct {
	client   *lunahttp.Client
	basePath string
}

// List retrieves all buckets
func (r *BucketsResource) List(ctx context.Context) (*BucketList, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   r.basePath,
	})
	if err != nil {
		return nil, err
	}

	var result BucketList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Upload uploads a file to a bucket
func (r *BucketsResource) Upload(ctx context.Context, bucketID string, file io.Reader, filename string, metadata map[string]string) (*FileObject, error) {
	// Create multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add metadata
	if len(metadata) > 0 {
		metaBytes, _ := json.Marshal(metadata)
		_ = writer.WriteField("metadata", string(metaBytes))
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%s/upload", r.basePath, bucketID),
		BodyReader:  body,
		ContentType: writer.FormDataContentType(),
	})
	if err != nil {
		return nil, err
	}

	var result FileObject
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// FilesResource provides access to file operations
type FilesResource struct {
	client   *lunahttp.Client
	basePath string
}

// GetDownloadURL retrieves the download URL for a file
func (r *FilesResource) GetDownloadURL(ctx context.Context, id string) (string, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s/download", r.basePath, id),
	})
	if err != nil {
		return "", err
	}

	// Assuming response is {"url": "..."}
	var result struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result.URL, nil
}

// StorageResource groups Storage service resources
type StorageResource struct {
	Buckets *BucketsResource
	Files   *FilesResource
}

// NewStorageResource creates a new Storage resource
func NewStorageResource(client *lunahttp.Client) *StorageResource {
	return &StorageResource{
		Buckets: &BucketsResource{
			client:   client,
			basePath: "/v1/storage/buckets",
		},
		Files: &FilesResource{
			client:   client,
			basePath: "/v1/storage/files",
		},
	}
}
