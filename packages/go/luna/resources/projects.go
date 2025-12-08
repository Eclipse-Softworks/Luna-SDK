package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

var projectIDPattern = regexp.MustCompile(`^prj_[a-zA-Z0-9]+$`)

// ProjectsResource provides access to project operations
type ProjectsResource struct {
	client   *lunahttp.Client
	basePath string
}

// NewProjectsResource creates a new projects resource
func NewProjectsResource(client *lunahttp.Client) *ProjectsResource {
	return &ProjectsResource{
		client:   client,
		basePath: "/v1/projects",
	}
}

// List retrieves all projects with pagination
func (r *ProjectsResource) List(ctx context.Context, params *ListParams) (*ProjectList, error) {
	query := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Cursor != "" {
			query.Set("cursor", params.Cursor)
		}
	}

	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   r.basePath,
		Query:  query,
	})
	if err != nil {
		return nil, err
	}

	var result ProjectList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Get retrieves a project by ID
func (r *ProjectsResource) Get(ctx context.Context, projectID string) (*Project, error) {
	if err := validateProjectID(projectID); err != nil {
		return nil, err
	}

	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", r.basePath, projectID),
	})
	if err != nil {
		return nil, err
	}

	var result Project
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Create creates a new project
func (r *ProjectsResource) Create(ctx context.Context, data ProjectCreate) (*Project, error) {
	if data.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "POST",
		Path:   r.basePath,
		Body:   data,
	})
	if err != nil {
		return nil, err
	}

	var result Project
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Update updates an existing project
func (r *ProjectsResource) Update(ctx context.Context, projectID string, data ProjectUpdate) (*Project, error) {
	if err := validateProjectID(projectID); err != nil {
		return nil, err
	}

	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "PATCH",
		Path:   fmt.Sprintf("%s/%s", r.basePath, projectID),
		Body:   data,
	})
	if err != nil {
		return nil, err
	}

	var result Project
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Delete deletes a project
func (r *ProjectsResource) Delete(ctx context.Context, projectID string) error {
	if err := validateProjectID(projectID); err != nil {
		return err
	}

	_, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%s", r.basePath, projectID),
	})
	return err
}

func validateProjectID(id string) error {
	if id == "" {
		return fmt.Errorf("project ID is required")
	}
	if !projectIDPattern.MatchString(id) {
		return fmt.Errorf("invalid project ID format, expected: prj_<id>")
	}
	return nil
}
