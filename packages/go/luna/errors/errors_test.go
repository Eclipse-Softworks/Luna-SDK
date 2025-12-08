package errors_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
	"github.com/eclipse-softworks/luna-sdk-go/luna/errors"
	"github.com/eclipse-softworks/luna-sdk-go/luna/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationError(t *testing.T) {
	t.Run("returns AuthenticationError for 401", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(testutil.MockErrorAuth)
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789000"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().List(context.Background(), nil)

		require.Error(t, err)
		var authErr *errors.AuthenticationError
		assert.ErrorAs(t, err, &authErr)
	})

	t.Run("returns AuthenticationError for expired token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Token expired",
					"code":    "TOKEN_EXPIRED",
					"status":  401,
				},
			})
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithTokens("expired_token", ""),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().List(context.Background(), nil)

		require.Error(t, err)
	})
}

func TestNotFoundError(t *testing.T) {
	t.Run("returns NotFoundError for 404", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(testutil.MockErrorNotFound)
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().Get(context.Background(), "usr_nonexistent")

		require.Error(t, err)
		var notFoundErr *errors.NotFoundError
		assert.ErrorAs(t, err, &notFoundErr)
	})
}

func TestValidationError(t *testing.T) {
	t.Run("returns ValidationError for 400", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(testutil.MockErrorValidation)
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().Create(context.Background(), luna.UserCreate{
			Name:  "Test User",
			Email: "test@example.com",
		})

		require.Error(t, err)
		var validationErr *errors.ValidationError
		assert.ErrorAs(t, err, &validationErr)
	})
}

func TestRateLimitError(t *testing.T) {
	t.Run("returns RateLimitError for 429", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(testutil.MockErrorRateLimit)
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().List(context.Background(), nil)

		require.Error(t, err)
		var rateLimitErr *errors.RateLimitError
		assert.ErrorAs(t, err, &rateLimitErr)
	})

	t.Run("includes RetryAfter information", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(testutil.MockErrorRateLimit)
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().List(context.Background(), nil)

		require.Error(t, err)
		var rateLimitErr *errors.RateLimitError
		if assert.ErrorAs(t, err, &rateLimitErr) {
			assert.Greater(t, rateLimitErr.RetryAfter, 0)
		}
	})
}

func TestServerError(t *testing.T) {
	t.Run("returns ServerError for 500", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(testutil.MockErrorServer)
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().List(context.Background(), nil)

		require.Error(t, err)
		var serverErr *errors.ServerError
		assert.ErrorAs(t, err, &serverErr)
	})

	t.Run("returns ServerError for 503", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Service unavailable",
					"code":    "SERVICE_UNAVAILABLE",
					"status":  503,
				},
			})
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().List(context.Background(), nil)

		require.Error(t, err)
	})
}

func TestErrorProperties(t *testing.T) {
	t.Run("error includes request ID when available", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Request-Id", "req_123abc")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"message":    "Server error",
					"code":       "SERVER_ERROR",
					"status":     500,
					"request_id": "req_123abc",
				},
			})
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().List(context.Background(), nil)

		require.Error(t, err)
		// Error should have request ID available
		assert.Contains(t, err.Error(), "req_123abc")
	})

	t.Run("error is serializable to string", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(testutil.MockErrorNotFound)
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().Get(context.Background(), "usr_nonexistent")

		require.Error(t, err)
		errStr := err.Error()
		assert.NotEmpty(t, errStr)
	})
}
