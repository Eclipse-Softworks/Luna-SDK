// Package resources provides API resource implementations for the Luna SDK.
package resources

import "time"

// ListParams holds common pagination parameters
type ListParams struct {
	Limit  int
	Cursor string
}

// User represents a user resource
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserCreate holds parameters for creating a user
type UserCreate struct {
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UserUpdate holds parameters for updating a user
type UserUpdate struct {
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UserList holds a paginated list of users
type UserList = ListResponse[User]

// Project represents a project resource
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectCreate holds parameters for creating a project
type ProjectCreate struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// ProjectUpdate holds parameters for updating a project
type ProjectUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ProjectList holds a paginated list of projects
type ProjectList = ListResponse[Project]
