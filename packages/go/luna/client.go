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
	User          = resources.User
	UserCreate    = resources.UserCreate
	UserUpdate    = resources.UserUpdate
	UserList      = resources.UserList
	Project       = resources.Project
	ProjectCreate = resources.ProjectCreate
	ProjectUpdate = resources.ProjectUpdate
	ProjectList   = resources.ProjectList
	ListParams    = resources.ListParams
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
}

// Client is the main Luna SDK client
type Client struct {
	config     *clientConfig
	httpClient *lunahttp.Client
	users      *resources.UsersResource
	projects   *resources.ProjectsResource
}

// NewClient creates a new Luna SDK client
func NewClient(opts ...Option) *Client {
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
	if config.apiKey == "" && config.accessToken == "" {
		panic("luna: either apiKey or accessToken must be provided")
	}

	// Set up logger
	logger := config.logger
	if logger == nil {
		logger = telemetry.NewConsoleLogger(config.logLevel)
	}

	// Set up auth provider
	var authProvider auth.Provider
	if config.apiKey != "" {
		authProvider = auth.NewAPIKeyAuth(config.apiKey)
	} else {
		authProvider = auth.NewTokenAuth(config.accessToken, config.refreshToken, config.tokenRefreshCallback)
	}

	// Create HTTP client
	httpClient := lunahttp.NewClient(lunahttp.ClientConfig{
		BaseURL:      config.baseURL,
		Timeout:      config.timeout,
		MaxRetries:   config.maxRetries,
		AuthProvider: authProvider,
		Logger:       logger,
	})

	client := &Client{
		config:     config,
		httpClient: httpClient,
	}

	// Initialize resources
	client.users = resources.NewUsersResource(httpClient)
	client.projects = resources.NewProjectsResource(httpClient)

	logger.Debug("LunaClient initialized", map[string]interface{}{
		"base_url":  config.baseURL,
		"auth_type": getAuthType(config),
	})

	return client
}

// Users returns the Users resource
func (c *Client) Users() *resources.UsersResource {
	return c.users
}

// Projects returns the Projects resource
func (c *Client) Projects() *resources.ProjectsResource {
	return c.projects
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

// WithLogLevel sets the log level
func WithLogLevel(level telemetry.LogLevel) Option {
	return func(c *clientConfig) {
		c.logLevel = level
	}
}

func getAuthType(c *clientConfig) string {
	if c.apiKey != "" {
		return "api_key"
	}
	return "token"
}
