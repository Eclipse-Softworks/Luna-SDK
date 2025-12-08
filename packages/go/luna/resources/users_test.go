package resources_test

import (
	"context"
	"testing"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
	"github.com/eclipse-softworks/luna-sdk-go/luna/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersResource_List(t *testing.T) {
	ms := testutil.NewMockServer()
	defer ms.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(ms.URL()),
	)
	require.NoError(t, err)

	t.Run("returns list of users", func(t *testing.T) {
		result, err := client.Users().List(context.Background(), nil)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, len(testutil.MockUsers))
	})

	t.Run("supports pagination with limit", func(t *testing.T) {
		result, err := client.Users().List(context.Background(), &luna.ListParams{
			Limit: 1,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("supports pagination with cursor", func(t *testing.T) {
		result, err := client.Users().List(context.Background(), &luna.ListParams{
			Cursor: "next_cursor",
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestUsersResource_Get(t *testing.T) {
	ms := testutil.NewMockServer()
	defer ms.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(ms.URL()),
	)
	require.NoError(t, err)

	t.Run("returns user by ID", func(t *testing.T) {
		user, err := client.Users().Get(context.Background(), "usr_123456789")

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "usr_123456789", user.ID)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		_, err = client.Users().Get(context.Background(), "usr_nonexistent")

		require.Error(t, err)
	})
}

func TestUsersResource_Create(t *testing.T) {
	ms := testutil.NewMockServer()
	defer ms.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(ms.URL()),
	)
	require.NoError(t, err)

	t.Run("creates new user", func(t *testing.T) {
		user, err := client.Users().Create(context.Background(), luna.UserCreate{
			Name:  "New User",
			Email: "newuser@example.com",
		})

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, user.ID)
	})
}

func TestUsersResource_Update(t *testing.T) {
	ms := testutil.NewMockServer()
	defer ms.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(ms.URL()),
	)
	require.NoError(t, err)

	t.Run("updates existing user", func(t *testing.T) {
		user, err := client.Users().Update(context.Background(), "usr_123456789", luna.UserUpdate{
			Name: stringPtr("Updated Name"),
		})

		require.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		_, err = client.Users().Update(context.Background(), "usr_nonexistent", luna.UserUpdate{
			Name: stringPtr("Test"),
		})

		require.Error(t, err)
	})
}

func TestUsersResource_Delete(t *testing.T) {
	ms := testutil.NewMockServer()
	defer ms.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(ms.URL()),
	)
	require.NoError(t, err)

	t.Run("deletes existing user", func(t *testing.T) {
		err := client.Users().Delete(context.Background(), "usr_123456789")

		require.NoError(t, err)
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		err := client.Users().Delete(context.Background(), "usr_nonexistent")

		require.Error(t, err)
	})
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
