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

// ResidenceLocation represents a residence location
type ResidenceLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Suburb    *string `json:"suburb,omitempty"`
	City      *string `json:"city,omitempty"`
}

// Residence represents a residence resource
type Residence struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Address     string  `json:"address"`
	Description *string `json:"description,omitempty"`

	// Filters & Attributes
	IsNSFASAccredited bool    `json:"is_nsfas_accredited"`
	MinPrice          float64 `json:"min_price"`
	MaxPrice          float64 `json:"max_price"`
	CurrencyCode      string  `json:"currency_code"`
	GenderPolicy      string  `json:"gender_policy"` // 'mixed' | 'male' | 'female'

	// Location & Relations
	Location  ResidenceLocation `json:"location"`
	CampusIDs []string          `json:"campus_ids"`

	// Social
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count"`

	Images    []string `json:"images"`
	Amenities []string `json:"amenities"`
}

// ResidenceSearch holds search/filter parameters for residences
type ResidenceSearch struct {
	ListParams
	Query     string  `json:"query,omitempty"`
	NSFAS     *bool   `json:"nsfas,omitempty"`
	MinPrice  float64 `json:"min_price,omitempty"`
	MaxPrice  float64 `json:"max_price,omitempty"`
	Gender    string  `json:"gender,omitempty"`
	CampusID  string  `json:"campus_id,omitempty"`
	Radius    float64 `json:"radius,omitempty"`
	MinRating float64 `json:"min_rating,omitempty"`
}

// ResidenceList holds a paginated list of residences
type ResidenceList = ListResponse[Residence]

// CampusLocation represents a campus location
type CampusLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Campus represents a campus resource
type Campus struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Location CampusLocation `json:"location"`
}

// CampusList holds a list of campuses
type CampusList = ListResponse[Campus]

// Group represents an identity group
type Group struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Permissions []string `json:"permissions"`
	MemberIDs   []string `json:"member_ids"`
}

// GroupCreate holds parameters for creating a group
type GroupCreate struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	MemberIDs   []string `json:"member_ids,omitempty"`
}

// GroupList holds a paginated list of groups
type GroupList = ListResponse[Group]

// Bucket represents a storage bucket
type Bucket struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

// BucketList holds a paginated list of buckets
type BucketList = ListResponse[Bucket]

// FileObject represents a file in storage
type FileObject struct {
	ID          string `json:"id"`
	BucketID    string `json:"bucket_id"`
	Key         string `json:"key"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

// CompletionRequest represents an AI completion request
type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature *float64  `json:"temperature,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponse represents an AI completion response
type CompletionResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

// Choice represents a completion choice
type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

// Workflow represents an automation workflow
type Workflow struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TriggerType string `json:"trigger_type"`
	IsActive    bool   `json:"is_active"`
}

// WorkflowList holds a paginated list of workflows
type WorkflowList = ListResponse[Workflow]

// WorkflowRun represents a workflow execution status
type WorkflowRun struct {
	ID         string `json:"id"`
	WorkflowID string `json:"workflow_id"`
	Status     string `json:"status"`
	StartedAt  string `json:"started_at"`
}
