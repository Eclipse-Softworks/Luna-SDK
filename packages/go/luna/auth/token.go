package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// TokenAuth implements OAuth token authentication with refresh
type TokenAuth struct {
	accessToken  string
	refreshToken string
	expiresAt    *time.Time
	mu           sync.RWMutex
}

// NewTokenAuth creates a new token authentication provider
func NewTokenAuth(accessToken, refreshToken string) *TokenAuth {
	if accessToken == "" {
		panic("auth: access token is required")
	}
	return &TokenAuth{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

// GetHeaders returns authorization headers with the access token
func (t *TokenAuth) GetHeaders() (map[string]string, error) {
	t.mu.RLock()
	token := t.accessToken
	t.mu.RUnlock()

	if t.NeedsRefresh() {
		if err := t.Refresh(); err != nil {
			return nil, err
		}
		t.mu.RLock()
		token = t.accessToken
		t.mu.RUnlock()
	}

	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}, nil
}

// NeedsRefresh returns true if the token is expiring soon
func (t *TokenAuth) NeedsRefresh() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.expiresAt == nil {
		return false
	}

	// Refresh if expiring within 5 minutes
	buffer := 5 * time.Minute
	return time.Now().Add(buffer).After(*t.expiresAt)
}

// Refresh refreshes the access token
func (t *TokenAuth) Refresh() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.refreshToken == "" {
		return fmt.Errorf("auth: no refresh token available")
	}

	body, _ := json.Marshal(map[string]string{
		"refresh_token": t.refreshToken,
	})

	resp, err := http.Post(
		"https://api.eclipse.dev/v1/auth/refresh",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("auth: refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth: refresh failed with status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("auth: failed to decode refresh response: %w", err)
	}

	t.accessToken = result.AccessToken
	t.refreshToken = result.RefreshToken
	expiresAt := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	t.expiresAt = &expiresAt

	return nil
}

// UpdateTokens manually updates the tokens
func (t *TokenAuth) UpdateTokens(accessToken, refreshToken string, expiresAt *time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.accessToken = accessToken
	t.refreshToken = refreshToken
	t.expiresAt = expiresAt
}

var _ Provider = (*TokenAuth)(nil)
