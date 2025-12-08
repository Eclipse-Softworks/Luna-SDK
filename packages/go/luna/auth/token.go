package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/eclipse-softworks/luna-sdk-go/luna/errors"
)

type TokenPair struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresAt    *time.Time `json:"expires_at"`
}

type TokenAuth struct {
	tokens    TokenPair
	onRefresh func(TokenPair) error
	mu        sync.RWMutex
}

func NewTokenAuth(access, refresh string, onRefresh func(TokenPair) error) *TokenAuth {
	return &TokenAuth{
		tokens:    TokenPair{AccessToken: access, RefreshToken: refresh},
		onRefresh: onRefresh,
	}
}

func (t *TokenAuth) GetHeaders() (map[string]string, error) {
	t.mu.RLock()
	token := t.tokens.AccessToken
	t.mu.RUnlock()

	if t.NeedsRefresh() {
		if err := t.Refresh(); err != nil {
			return nil, err
		}
		t.mu.RLock()
		token = t.tokens.AccessToken
		t.mu.RUnlock()
	}

	return map[string]string{"Authorization": "Bearer " + token}, nil
}

func (t *TokenAuth) NeedsRefresh() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.tokens.ExpiresAt == nil {
		return false
	}
	return time.Now().Add(5 * time.Minute).After(*t.tokens.ExpiresAt)
}

func (t *TokenAuth) Refresh() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Double check
	if t.tokens.RefreshToken == "" {
		return &errors.AuthenticationError{
			BaseError: errors.New(errors.CodeAuthInvalidKey, "no refresh token available", 401, "local", nil),
		}
	}

	// Perform HTTP request (example)
	body, _ := json.Marshal(map[string]string{"refresh_token": t.tokens.RefreshToken})
	resp, err := http.Post("https://api.eclipse.dev/v1/auth/refresh", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return &errors.AuthenticationError{
			BaseError: errors.New(errors.CodeAuthTokenExpired, fmt.Sprintf("refresh failed: %d", resp.StatusCode), 401, resp.Header.Get("X-Request-Id"), nil),
		}
	}

	var res struct {
		Access  string `json:"access_token"`
		Refresh string `json:"refresh_token"`
		Expires int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	expiry := time.Now().Add(time.Duration(res.Expires) * time.Second)
	t.tokens = TokenPair{
		AccessToken:  res.Access,
		RefreshToken: res.Refresh,
		ExpiresAt:    &expiry,
	}

	// Trigger callback for persistence
	if t.onRefresh != nil {
		go t.onRefresh(t.tokens) // Run async or sync depending on strictness requirements
	}

	return nil
}
