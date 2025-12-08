// Package luna provides the official Go SDK for the Eclipse Softworks Platform API.
//
// Quick Start:
//
//	client := luna.NewClient(luna.WithAPIKey(os.Getenv("LUNA_API_KEY")))
//
//	users, err := client.Users().List(ctx, &luna.ListParams{Limit: 10})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	user, err := client.Users().Get(ctx, "usr_123")
package luna

import (
	"fmt"

	"github.com/eclipse-softworks/luna-sdk-go/luna/auth"
	"github.com/eclipse-softworks/luna-sdk-go/luna/errors"
	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
	"github.com/eclipse-softworks/luna-sdk-go/luna/resources"
	"github.com/eclipse-softworks/luna-sdk-go/luna/telemetry"
)

// Version of the SDK
const Version = "1.0.0"

// Re-export commonly used types
type (
	// Error types
	Error               = errors.Error
	AuthenticationError = errors.AuthenticationError
	AuthorizationError  = errors.AuthorizationError
	ValidationError     = errors.ValidationError
	RateLimitError      = errors.RateLimitError
	NetworkError        = errors.NetworkError
	NotFoundError       = errors.NotFoundError
	ConflictError       = errors.ConflictError
	ServerError         = errors.ServerError

	// Resource types
	User               = resources.User
	UserCreate         = resources.UserCreate
	UserUpdate         = resources.UserUpdate
	UserList           = resources.UserList
	Project            = resources.Project
	ProjectCreate      = resources.ProjectCreate
	ProjectUpdate      = resources.ProjectUpdate
	ProjectList        = resources.ProjectList
	ListParams         = resources.ListParams
	Residence          = resources.Residence
	ResidenceList      = resources.ResidenceList
	Campus             = resources.Campus
	CampusList         = resources.CampusList
	Group              = resources.Group
	GroupCreate        = resources.GroupCreate
	GroupList          = resources.GroupList
	Bucket             = resources.Bucket
	BucketList         = resources.BucketList
	FileObject         = resources.FileObject
	CompletionRequest  = resources.CompletionRequest
	CompletionResponse = resources.CompletionResponse
	Message            = resources.Message
	Choice             = resources.Choice
	Workflow           = resources.Workflow
	WorkflowList       = resources.WorkflowList
	WorkflowRun        = resources.WorkflowRun
)

// Client configuration options
type (
	Option       = func(*Config)
	Config       = clientConfig
	Logger       = telemetry.Logger
	LogLevel     = telemetry.LogLevel
	AuthProvider = auth.Provider
)

// clientConfig holds client configuration
type clientConfig struct {
	apiKey               string
	accessToken          string
	refreshToken         string
	baseURL              string
	timeout              int
	maxRetries           int
	logger               telemetry.Logger
	logLevel             telemetry.LogLevel
	tokenRefreshCallback func(auth.TokenPair) error
	httpClient           *lunahttp.Client
}

// Client is the main Luna SDK client
type Client struct {
	config     *clientConfig
	httpClient *lunahttp.Client
	users      *resources.UsersResource
	projects   *resources.ProjectsResource
	resMate    *resources.ResMateResource
	identity   *resources.IdentityResource
	storage    *resources.StorageResource
	ai         *resources.AiResource
	automation *resources.AutomationResource
}

// NewClient creates a new Luna SDK client
func NewClient(opts ...Option) (*Client, error) {
	config := &clientConfig{
		baseURL:    "https://api.eclipse.dev",
		timeout:    30000,
		maxRetries: 3,
		logLevel:   telemetry.LogLevelInfo,
	}

	for _, opt := range opts {
		opt(config)
	}

	// Validate config
	if config.baseURL == "" {
		return nil, fmt.Errorf("luna: base URL cannot be empty")
	}

	// Validate auth
	if config.apiKey == "" && config.accessToken == "" {
		return nil, fmt.Errorf("luna: either apiKey or accessToken must be provided")
	}

	// Set up logger
	logger := config.logger
	if logger == nil {
		logger = telemetry.NewConsoleLogger(config.logLevel)
	}

	// Set up auth provider
	var authProvider auth.Provider
	if config.apiKey != "" {
		var err error
		authProvider, err = auth.NewAPIKeyAuth(config.apiKey)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		authProvider, err = auth.NewTokenAuth(config.accessToken, config.refreshToken, config.tokenRefreshCallback)
		if err != nil {
			return nil, err
		}
	}

	// Create HTTP client
	var httpClient *lunahttp.Client
	if config.httpClient != nil {
		httpClient = config.httpClient
	} else {
		// Allow for custom HTTP client if we add that option later, but for now use standard
		httpClient = lunahttp.NewClient(lunahttp.ClientConfig{
			BaseURL:      config.baseURL,
			Timeout:      config.timeout,
			MaxRetries:   config.maxRetries,
			AuthProvider: authProvider,
			Logger:       logger,
		})
	}

	client := &Client{
		config:     config,
		httpClient: httpClient,
	}

	// Initialize resources
	client.users = resources.NewUsersResource(httpClient)
	client.projects = resources.NewProjectsResource(httpClient)
	client.resMate = resources.NewResMateResource(httpClient)
	client.identity = resources.NewIdentityResource(httpClient)
	client.storage = resources.NewStorageResource(httpClient)
	client.ai = resources.NewAiResource(httpClient)
	client.automation = resources.NewAutomationResource(httpClient)

	logger.Debug("LunaClient initialized", map[string]interface{}{
		"base_url":  config.baseURL,
		"auth_type": getAuthType(config),
	})

	return client, nil
}

// Users returns the Users resource
func (c *Client) Users() *resources.UsersResource {
	return c.users
}

// Projects returns the Projects resource
func (c *Client) Projects() *resources.ProjectsResource {
	return c.projects
}

// ResMate returns the ResMate resource
func (c *Client) ResMate() *resources.ResMateResource {
	return c.resMate
}

// Identity returns the Identity resource
func (c *Client) Identity() *resources.IdentityResource {
	return c.identity
}

// Storage returns the Storage resource
func (c *Client) Storage() *resources.StorageResource {
	return c.storage
}

// AI returns the AI resource
func (c *Client) AI() *resources.AiResource {
	return c.ai
}

// Automation returns the Automation resource
func (c *Client) Automation() *resources.AutomationResource {
	return c.automation
}

// WithAPIKey sets the API key for authentication
func WithAPIKey(apiKey string) Option {
	return func(c *clientConfig) {
		c.apiKey = apiKey
	}
}

// WithTokens sets OAuth tokens for authentication
func WithTokens(accessToken, refreshToken string) Option {
	return func(c *clientConfig) {
		c.accessToken = accessToken
		c.refreshToken = refreshToken
	}
}

// WithTokenRefreshCallback sets the callback for token refresh events
func WithTokenRefreshCallback(callback func(auth.TokenPair) error) Option {
	return func(c *clientConfig) {
		c.tokenRefreshCallback = callback
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) Option {
	return func(c *clientConfig) {
		c.baseURL = baseURL
	}
}

// WithTimeout sets the request timeout in milliseconds
func WithTimeout(timeout int) Option {
	return func(c *clientConfig) {
		c.timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retry attempts
func WithMaxRetries(maxRetries int) Option {
	return func(c *clientConfig) {
		c.maxRetries = maxRetries
	}
}

// WithLogger sets a custom logger
func WithLogger(logger telemetry.Logger) Option {
	return func(c *clientConfig) {
		c.logger = logger
	}
}

// WithHTTPClient allows providing a custom HTTP client
func WithHTTPClient(client *lunahttp.Client) Option {
	return func(c *clientConfig) {
		c.httpClient = client
	}
}

func getAuthType(c *clientConfig) string {
	if c.apiKey != "" {
		return "api_key"
	}
	return "token"
}
