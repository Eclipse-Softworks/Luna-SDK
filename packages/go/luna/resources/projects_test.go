package resources_test

import (
	"context"
	"testing"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
	"github.com/eclipse-softworks/luna-sdk-go/luna/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectsResource_List(t *testing.T) {
	ms := testutil.NewMockServer()
	defer ms.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(ms.URL()),
	)
	require.NoError(t, err)

	t.Run("returns list of projects", func(t *testing.T) {
		result, err := client.Projects().List(context.Background(), nil)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Data, len(testutil.MockProjects))
	})

	t.Run("supports pagination parameters", func(t *testing.T) {
		result, err := client.Projects().List(context.Background(), &luna.ListParams{
			Limit: 5,
		})

		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestProjectsResource_Get(t *testing.T) {
	ms := testutil.NewMockServer()
	defer ms.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(ms.URL()),
	)
	require.NoError(t, err)

	t.Run("returns project by ID", func(t *testing.T) {
		project, err := client.Projects().Get(context.Background(), "prj_123456789")

		require.NoError(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "prj_123456789", project.ID)
	})

	t.Run("returns error for non-existent project", func(t *testing.T) {
		_, err := client.Projects().Get(context.Background(), "prj_nonexistent")

		require.Error(t, err)
	})
}

func TestProjectsResource_Create(t *testing.T) {
	ms := testutil.NewMockServer()
	defer ms.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(ms.URL()),
	)
	require.NoError(t, err)

	t.Run("creates new project", func(t *testing.T) {
		description := "A new project"
		project, err := client.Projects().Create(context.Background(), luna.ProjectCreate{
			Name:        "New Project",
			Description: &description,
		})

		require.NoError(t, err)
		assert.NotNil(t, project)
		assert.NotEmpty(t, project.ID)
	})

	t.Run("creates project with minimal fields", func(t *testing.T) {
		project, err := client.Projects().Create(context.Background(), luna.ProjectCreate{
			Name: "Minimal Project",
		})

		require.NoError(t, err)
		assert.NotNil(t, project)
	})
}

func TestProjectsResource_Delete(t *testing.T) {
	ms := testutil.NewMockServer()
	defer ms.Close()

	client, err := luna.NewClient(
		luna.WithAPIKey("lk_test_12345678901234567890123456789012"),
		luna.WithBaseURL(ms.URL()),
	)
	require.NoError(t, err)

	t.Run("deletes existing project", func(t *testing.T) {
		err := client.Projects().Delete(context.Background(), "prj_123456789")

		require.NoError(t, err)
	})

	t.Run("returns error for non-existent project", func(t *testing.T) {
		err := client.Projects().Delete(context.Background(), "prj_nonexistent")

		require.Error(t, err)
	})
}
