package auth

import (
	"fmt"
	"regexp"
)

// APIKeyAuth implements API key authentication
type APIKeyAuth struct {
	apiKey string
}

var apiKeyPattern = regexp.MustCompile(`^lk_(live|test|dev)_[a-zA-Z0-9]{32}$`)

// NewAPIKeyAuth creates a new API key authentication provider
func NewAPIKeyAuth(apiKey string) (*APIKeyAuth, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("auth: API key is required")
	}
	if !apiKeyPattern.MatchString(apiKey) {
		return nil, fmt.Errorf("auth: invalid API key format, expected: lk_<env>_<key>")
	}
	return &APIKeyAuth{apiKey: apiKey}, nil
}

// GetHeaders returns authorization headers with the API key
func (a *APIKeyAuth) GetHeaders() (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + a.apiKey,
	}, nil
}

// NeedsRefresh returns false as API keys don't expire
func (a *APIKeyAuth) NeedsRefresh() bool {
	return false
}

// Refresh is a no-op for API keys
func (a *APIKeyAuth) Refresh() error {
	return nil
}

var _ Provider = (*APIKeyAuth)(nil)
