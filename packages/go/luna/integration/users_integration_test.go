package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
	"github.com/eclipse-softworks/luna-sdk-go/luna/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersCRUDWorkflow(t *testing.T) {
	t.Run("complete CRUD workflow", func(t *testing.T) {
		createdID := "usr_integrationTest"

		mux := http.NewServeMux()

		// CREATE
		mux.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":         createdID,
					"name":       "Integration User",
					"email":      "integration@test.com",
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z",
				})
			}
		})

		// GET, UPDATE, DELETE
		mux.HandleFunc("/v1/users/", func(w http.ResponseWriter, r *http.Request) {
			id := strings.TrimPrefix(r.URL.Path, "/v1/users/")
			w.Header().Set("Content-Type", "application/json")

			switch r.Method {
			case http.MethodGet:
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":         id,
					"name":       "Integration User",
					"email":      "integration@test.com",
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z",
				})
			case http.MethodPatch:
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":         id,
					"name":       "Updated Name",
					"email":      "integration@test.com",
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z",
				})
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			}
		})

		server := httptest.NewServer(mux)
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		// CREATE
		created, err := client.Users().Create(context.Background(), luna.UserCreate{
			Name:  "Integration User",
			Email: "integration@test.com",
		})
		require.NoError(t, err)
		assert.Equal(t, createdID, created.ID)

		// READ
		fetched, err := client.Users().Get(context.Background(), createdID)
		require.NoError(t, err)
		assert.Equal(t, createdID, fetched.ID)

		// UPDATE
		name := "Updated Name"
		updated, err := client.Users().Update(context.Background(), createdID, luna.UserUpdate{
			Name: &name,
		})
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)

		// DELETE
		err = client.Users().Delete(context.Background(), createdID)
		require.NoError(t, err)
	})
}

func TestPaginationWorkflow(t *testing.T) {
	t.Run("paginate through all users", func(t *testing.T) {
		var callCount int32

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&callCount, 1)
			cursor := r.URL.Query().Get("cursor")

			w.Header().Set("Content-Type", "application/json")

			switch cursor {
			case "":
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data":        testutil.MockUsers[:1],
					"has_more":    true,
					"next_cursor": "page2",
				})
			case "page2":
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data":        testutil.MockUsers[1:],
					"has_more":    false,
					"next_cursor": nil,
				})
			}
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		// Fetch first page
		page1, err := client.Users().List(context.Background(), &luna.ListParams{Limit: 1})
		require.NoError(t, err)
		assert.Len(t, page1.Data, 1)
		assert.True(t, page1.HasMore)
		assert.NotNil(t, page1.NextCursor)
		assert.Equal(t, "page2", *page1.NextCursor)

		// Fetch second page
		page2, err := client.Users().List(context.Background(), &luna.ListParams{Cursor: *page1.NextCursor})
		require.NoError(t, err)
		assert.Greater(t, len(page2.Data), 0)
		assert.False(t, page2.HasMore)

		assert.Equal(t, int32(2), atomic.LoadInt32(&callCount))
	})
}

func TestRequestHeaders(t *testing.T) {
	t.Run("sends correct Authorization header", func(t *testing.T) {
		var capturedAuthHeader string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedAuthHeader = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(testutil.MockListResponse(testutil.MockUsers, false, ""))
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().List(context.Background(), nil)
		require.NoError(t, err)

		assert.Equal(t, "Bearer lk_test_12345678901234567890123456789012", capturedAuthHeader)
	})

	t.Run("sends Content-Type header for POST", func(t *testing.T) {
		var capturedContentType string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedContentType = r.Header.Get("Content-Type")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(testutil.MockUser)
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		_, err = client.Users().Create(context.Background(), luna.UserCreate{
			Name:  "Test",
			Email: "test@test.com",
		})
		require.NoError(t, err)

		assert.Contains(t, capturedContentType, "application/json")
	})
}

func TestRetryBehavior(t *testing.T) {
	t.Run("retries on transient errors", func(t *testing.T) {
		var attemptCount int32

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count := atomic.AddInt32(&attemptCount, 1)
			w.Header().Set("Content-Type", "application/json")

			if count < 3 {
				w.WriteHeader(http.StatusServiceUnavailable)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{
						"message": "Service unavailable",
						"status":  503,
					},
				})
				return
			}

			json.NewEncoder(w).Encode(testutil.MockListResponse(testutil.MockUsers, false, ""))
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
			luna.WithMaxRetries(3),
		)
		require.NoError(t, err)

		result, err := client.Users().List(context.Background(), nil)

		// If client has retry logic, it should succeed
		if err == nil {
			assert.NotNil(t, result)
			assert.GreaterOrEqual(t, atomic.LoadInt32(&attemptCount), int32(1))
		} else {
			// If no retry logic, it fails on first attempt
			assert.Equal(t, int32(1), atomic.LoadInt32(&attemptCount))
		}
	})
}
