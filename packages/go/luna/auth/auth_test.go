package auth_test

import (
	"testing"

	"github.com/eclipse-softworks/luna-sdk-go/luna/auth"
)

const validAPIKey = "lk_test_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestAPIKeyAuth(t *testing.T) {
	t.Run("creates auth with valid key", func(t *testing.T) {
		a, err := auth.NewAPIKeyAuth(validAPIKey)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a == nil {
			t.Fatal("expected auth to be created")
		}
	})

	t.Run("returns error on empty key", func(t *testing.T) {
		_, err := auth.NewAPIKeyAuth("")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns error on invalid format", func(t *testing.T) {
		_, err := auth.NewAPIKeyAuth("invalid-key")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns correct headers", func(t *testing.T) {
		a, err := auth.NewAPIKeyAuth(validAPIKey)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		headers, err := a.GetHeaders()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if headers["Authorization"] != "Bearer "+validAPIKey {
			t.Errorf("expected Authorization=Bearer %s, got %s", validAPIKey, headers["Authorization"])
		}
	})

	t.Run("does not need refresh", func(t *testing.T) {
		a, err := auth.NewAPIKeyAuth(validAPIKey)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a.NeedsRefresh() {
			t.Error("expected NeedsRefresh to be false")
		}
	})
}

func TestTokenAuth(t *testing.T) {
	t.Run("creates auth with access token", func(t *testing.T) {
		a, err := auth.NewTokenAuth("access-token", "refresh-token", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a == nil {
			t.Fatal("expected auth to be created")
		}
	})

	t.Run("returns error on empty access token", func(t *testing.T) {
		_, err := auth.NewTokenAuth("", "refresh-token", nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("returns correct headers", func(t *testing.T) {
		a, err := auth.NewTokenAuth("access-token", "refresh-token", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
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
		a, err := auth.NewTokenAuth("access-token", "refresh-token", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a.NeedsRefresh() {
			t.Error("expected NeedsRefresh to be false")
		}
	})
}
