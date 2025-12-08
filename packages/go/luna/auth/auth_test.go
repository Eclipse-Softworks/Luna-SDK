package auth_test

import (
	"testing"

	"github.com/eclipse-softworks/luna-sdk-go/luna/auth"
)

const validAPIKey = "lk_test_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestAPIKeyAuth(t *testing.T) {
	t.Run("creates auth with valid key", func(t *testing.T) {
		a := auth.NewAPIKeyAuth(validAPIKey)
		if a == nil {
			t.Fatal("expected auth to be created")
		}
	})

	t.Run("panics on empty key", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		auth.NewAPIKeyAuth("")
	})

	t.Run("panics on invalid format", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		auth.NewAPIKeyAuth("invalid-key")
	})

	t.Run("returns correct headers", func(t *testing.T) {
		a := auth.NewAPIKeyAuth(validAPIKey)
		headers, err := a.GetHeaders()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := map[string]string{"X-Luna-Api-Key": validAPIKey}
		for k, v := range expected {
			if headers[k] != v {
				t.Errorf("expected %s=%s, got %s", k, v, headers[k])
			}
		}
	})

	t.Run("does not need refresh", func(t *testing.T) {
		a := auth.NewAPIKeyAuth(validAPIKey)
		if a.NeedsRefresh() {
			t.Error("expected NeedsRefresh to be false")
		}
	})
}

func TestTokenAuth(t *testing.T) {
	t.Run("creates auth with access token", func(t *testing.T) {
		a := auth.NewTokenAuth("access-token", "refresh-token", nil)
		if a == nil {
			t.Fatal("expected auth to be created")
		}
	})

	t.Run("panics on empty access token", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		auth.NewTokenAuth("", "refresh-token", nil)
	})

	t.Run("returns correct headers", func(t *testing.T) {
		a := auth.NewTokenAuth("access-token", "refresh-token", nil)
		headers, err := a.GetHeaders()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "Bearer access-token"
		if headers["Authorization"] != expected {
			t.Errorf("expected Authorization=%s, got %s", expected, headers["Authorization"])
		}
	})

	t.Run("does not need refresh without expiry", func(t *testing.T) {
		a := auth.NewTokenAuth("access-token", "refresh-token", nil)
		if a.NeedsRefresh() {
			t.Error("expected NeedsRefresh to be false")
		}
	})
}
