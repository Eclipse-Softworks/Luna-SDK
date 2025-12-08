package luna_test

import (
	"testing"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
	"github.com/eclipse-softworks/luna-sdk-go/luna/errors"
)

const validAPIKey = "lk_test_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestNewClient(t *testing.T) {
	t.Run("creates client with API key", func(t *testing.T) {
		client := luna.NewClient(luna.WithAPIKey(validAPIKey))
		if client == nil {
			t.Fatal("expected client to be created")
		}
	})

	t.Run("creates client with tokens", func(t *testing.T) {
		client := luna.NewClient(luna.WithTokens("access-token", "refresh-token"))
		if client == nil {
			t.Fatal("expected client to be created")
		}
	})

	t.Run("panics without auth", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		luna.NewClient()
	})

	t.Run("accepts custom base URL", func(t *testing.T) {
		client := luna.NewClient(
			luna.WithAPIKey(validAPIKey),
			luna.WithBaseURL("https://api.staging.eclipse.dev"),
		)
		if client == nil {
			t.Fatal("expected client to be created")
		}
	})

	t.Run("accepts custom timeout", func(t *testing.T) {
		client := luna.NewClient(
			luna.WithAPIKey(validAPIKey),
			luna.WithTimeout(60000),
		)
		if client == nil {
			t.Fatal("expected client to be created")
		}
	})

	t.Run("accepts custom max retries", func(t *testing.T) {
		client := luna.NewClient(
			luna.WithAPIKey(validAPIKey),
			luna.WithMaxRetries(5),
		)
		if client == nil {
			t.Fatal("expected client to be created")
		}
	})
}

func TestClientResources(t *testing.T) {
	client := luna.NewClient(luna.WithAPIKey(validAPIKey))

	t.Run("exposes users resource", func(t *testing.T) {
		if client.Users() == nil {
			t.Fatal("expected users resource")
		}
	})

	t.Run("exposes projects resource", func(t *testing.T) {
		if client.Projects() == nil {
			t.Fatal("expected projects resource")
		}
	})
}

func TestError(t *testing.T) {
	t.Run("creates error with properties", func(t *testing.T) {
		err := errors.New(
			"LUNA_ERR_TEST",
			"Test error",
			500,
			"req_123",
			nil,
		)

		if err.Code != "LUNA_ERR_TEST" {
			t.Errorf("expected code LUNA_ERR_TEST, got %s", err.Code)
		}
		if err.Message != "Test error" {
			t.Errorf("expected message 'Test error', got %s", err.Message)
		}
		if err.Status != 500 {
			t.Errorf("expected status 500, got %d", err.Status)
		}
		if err.RequestID != "req_123" {
			t.Errorf("expected request ID req_123, got %s", err.RequestID)
		}
	})

	t.Run("generates docs URL", func(t *testing.T) {
		err := errors.New("LUNA_ERR_TEST", "Test", 500, "req_123", nil)
		expected := "https://docs.eclipse.dev/luna/errors#LUNA_ERR_TEST"
		if err.DocsURL() != expected {
			t.Errorf("expected %s, got %s", expected, err.DocsURL())
		}
	})

	t.Run("identifies retryable errors", func(t *testing.T) {
		serverErr := errors.New(errors.CodeServerInternal, "Internal", 500, "req_123", nil)
		authErr := errors.New(errors.CodeAuthInvalidKey, "Invalid", 401, "req_123", nil)

		if !serverErr.Retryable() {
			t.Error("expected server error to be retryable")
		}
		if authErr.Retryable() {
			t.Error("expected auth error to not be retryable")
		}
	})
}
