package resources

import (
	"context"
	"encoding/json"
	"fmt"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// GroupsResource provides access to group operations
type GroupsResource struct {
	client   *lunahttp.Client
	basePath string
}

// List retrieves all groups
func (r *GroupsResource) List(ctx context.Context) (*GroupList, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   r.basePath,
	})
	if err != nil {
		return nil, err
	}

	var result GroupList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Get retrieves a group by ID
func (r *GroupsResource) Get(ctx context.Context, id string) (*Group, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", r.basePath, id),
	})
	if err != nil {
		return nil, err
	}

	var result Group
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Create creates a new group
func (r *GroupsResource) Create(ctx context.Context, params *GroupCreate) (*Group, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "POST",
		Path:   r.basePath,
		Body:   params,
	})
	if err != nil {
		return nil, err
	}

	var result Group
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// IdentityResource groups Identity service resources
type IdentityResource struct {
	Groups *GroupsResource
}

// NewIdentityResource creates a new Identity resource
func NewIdentityResource(client *lunahttp.Client) *IdentityResource {
	return &IdentityResource{
		Groups: &GroupsResource{
			client:   client,
			basePath: "/v1/identity/groups",
		},
	}
}
