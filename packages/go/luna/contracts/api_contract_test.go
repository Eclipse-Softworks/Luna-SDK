package contracts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
	"github.com/eclipse-softworks/luna-sdk-go/luna/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testutil.MockUser)
	}))
	defer server.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	t.Run("user has required fields", func(t *testing.T) {
		user, err := client.Users().Get(context.Background(), "usr_123456789")

		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Name)
		assert.NotEmpty(t, user.Email)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
	})

	t.Run("user ID has correct prefix", func(t *testing.T) {
		user, err := client.Users().Get(context.Background(), "usr_123456789")

		require.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile(`^usr_`), user.ID)
	})

	t.Run("user email is valid format", func(t *testing.T) {
		user, err := client.Users().Get(context.Background(), "usr_123456789")

		require.NoError(t, err)
		emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
		assert.Regexp(t, emailRegex, user.Email)
	})
}

func TestProjectContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testutil.MockProject)
	}))
	defer server.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	t.Run("project has required fields", func(t *testing.T) {
		project, err := client.Projects().Get(context.Background(), "prj_123456789")

		require.NoError(t, err)
		assert.NotEmpty(t, project.ID)
		assert.NotEmpty(t, project.Name)
		assert.NotNil(t, project.CreatedAt)
		assert.NotNil(t, project.UpdatedAt)
	})

	t.Run("project ID has correct prefix", func(t *testing.T) {
		project, err := client.Projects().Get(context.Background(), "prj_123456789")

		require.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile(`^prj_`), project.ID)
	})

	t.Run("project description is optional", func(t *testing.T) {
		// Create server that returns project without description
		noDescServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			project := testutil.MockProject
			project["description"] = nil
			json.NewEncoder(w).Encode(project)
		}))
		defer noDescServer.Close()

		noDescClient, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(noDescServer.URL),
		)
		require.NoError(t, err)

		project, err := noDescClient.Projects().Get(context.Background(), "prj_123456789")

		require.NoError(t, err)
		assert.NotEmpty(t, project.ID)
		// Description can be nil or empty
	})
}

func TestBucketContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testutil.MockListResponse(testutil.MockBuckets, false, ""))
	}))
	defer server.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	t.Run("bucket list returns buckets", func(t *testing.T) {
		buckets, err := client.Storage().Buckets.List(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, buckets)
		assert.Greater(t, len(buckets.Data), 0)
	})

	t.Run("bucket has required fields", func(t *testing.T) {
		buckets, err := client.Storage().Buckets.List(context.Background())

		require.NoError(t, err)
		if len(buckets.Data) > 0 {
			bucket := buckets.Data[0]
			assert.NotEmpty(t, bucket.ID)
			assert.NotEmpty(t, bucket.Name)
			assert.Regexp(t, regexp.MustCompile(`^bkt_`), bucket.ID)
		}
	})
}

func TestListResponseContract(t *testing.T) {
	t.Run("list response has correct structure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(testutil.MockListResponse(testutil.MockUsers, false, ""))
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		result, err := client.Users().List(context.Background(), nil)

		require.NoError(t, err)
		assert.NotNil(t, result.Data)
		// HasMore should be a boolean (checked by type system)
	})

	t.Run("list response includes next_cursor when has_more is true", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(testutil.MockListResponse(testutil.MockUsers, true, "cursor_abc123"))
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		result, err := client.Users().List(context.Background(), nil)

		require.NoError(t, err)
		assert.True(t, result.HasMore)
		assert.NotEmpty(t, result.NextCursor)
	})
}

func TestErrorResponseContract(t *testing.T) {
	t.Run("error response has correct structure", func(t *testing.T) {
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
		// Error should have message
		assert.NotEmpty(t, err.Error())
	})

	t.Run("validation error includes details", func(t *testing.T) {
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
		// Error message should contain validation info
		assert.Contains(t, err.Error(), "Validation")
	})
}

func TestTimestampFormats(t *testing.T) {
	t.Run("timestamps are valid", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(testutil.MockUser)
		}))
		defer server.Close()

		client, err := luna.NewClient(
			luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
			luna.WithBaseURL(server.URL),
		)
		require.NoError(t, err)

		user, err := client.Users().Get(context.Background(), "usr_123456789")

		require.NoError(t, err)
		// CreatedAt should be a valid time
		assert.False(t, user.CreatedAt.IsZero())
	})
}
