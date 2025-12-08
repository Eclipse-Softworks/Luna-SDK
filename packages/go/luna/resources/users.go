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

var userIDPattern = regexp.MustCompile(`^usr_[a-zA-Z0-9]+$`)

// UsersResource provides access to user operations
type UsersResource struct {
	client   *lunahttp.Client
	basePath string
}

// NewUsersResource creates a new users resource
func NewUsersResource(client *lunahttp.Client) *UsersResource {
	return &UsersResource{
		client:   client,
		basePath: "/v1/users",
	}
}

// List retrieves all users with pagination
func (r *UsersResource) List(ctx context.Context, params *ListParams) (*ListResponse[User], error) {
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

	var result UserList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Iterate returns a paginator for iterating over users
func (r *UsersResource) Iterate(ctx context.Context, params *ListParams) *Paginator[User] {
	return NewPaginator(ctx, func(ctx context.Context, cursor string) (*ListResponse[User], error) {
		p := params
		if p == nil {
			p = &ListParams{}
		}
		p.Cursor = cursor
		return r.List(ctx, p)
	})
}

// Get retrieves a user by ID
func (r *UsersResource) Get(ctx context.Context, userID string) (*User, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", r.basePath, userID),
	})
	if err != nil {
		return nil, err
	}

	var result User
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Create creates a new user
func (r *UsersResource) Create(ctx context.Context, data UserCreate) (*User, error) {
	if data.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
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

	var result User
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Update updates an existing user
func (r *UsersResource) Update(ctx context.Context, userID string, data UserUpdate) (*User, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "PATCH",
		Path:   fmt.Sprintf("%s/%s", r.basePath, userID),
		Body:   data,
	})
	if err != nil {
		return nil, err
	}

	var result User
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Delete deletes a user
func (r *UsersResource) Delete(ctx context.Context, userID string) error {
	if err := validateUserID(userID); err != nil {
		return err
	}

	_, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%s", r.basePath, userID),
	})
	return err
}

func validateUserID(id string) error {
	if id == "" {
		return fmt.Errorf("user ID is required")
	}
	if !userIDPattern.MatchString(id) {
		return fmt.Errorf("invalid user ID format, expected: usr_<id>")
	}
	return nil
}
