// Package main demonstrates basic usage of the Luna SDK for Go.
//
// Run with: go run basic_usage.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
)

func main() {
	// Initialize the client with API key authentication
	client := luna.NewClient(
		luna.WithAPIKey(os.Getenv("LUNA_API_KEY")),
		luna.WithTimeout(30000),
		luna.WithMaxRetries(3),
	)

	ctx := context.Background()

	// Run all examples
	if err := userManagementExample(ctx, client); err != nil {
		log.Printf("User management example failed: %v", err)
	}

	if err := projectManagementExample(ctx, client); err != nil {
		log.Printf("Project management example failed: %v", err)
	}

	if err := paginationExample(ctx, client); err != nil {
		log.Printf("Pagination example failed: %v", err)
	}

	if err := errorHandlingExample(ctx, client); err != nil {
		log.Printf("Error handling example failed: %v", err)
	}
}

// ============================================
// Example 1: User Management
// ============================================

func userManagementExample(ctx context.Context, client *luna.Client) error {
	fmt.Println("=== User Management ===\n")

	// List users with pagination
	userList, err := client.Users().List(ctx, &luna.ListParams{Limit: 10})
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}
	fmt.Printf("Found %d users\n", len(userList.Data))

	// Create a new user
	newUser, err := client.Users().Create(ctx, luna.UserCreate{
		Email: "jane.doe@example.com",
		Name:  "Jane Doe",
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Printf("Created user: %s\n", newUser.ID)

	// Get user details
	user, err := client.Users().Get(ctx, newUser.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	fmt.Printf("User name: %s, Email: %s\n", user.Name, user.Email)

	// Update the user
	avatarURL := "https://example.com/avatar.jpg"
	updatedUser, err := client.Users().Update(ctx, newUser.ID, luna.UserUpdate{
		Name:      stringPtr("Jane M. Doe"),
		AvatarURL: &avatarURL,
	})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	fmt.Printf("Updated user name: %s\n", updatedUser.Name)

	// Delete the user
	if err := client.Users().Delete(ctx, newUser.ID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	fmt.Println("User deleted")

	return nil
}

// ============================================
// Example 2: Project Management
// ============================================

func projectManagementExample(ctx context.Context, client *luna.Client) error {
	fmt.Println("\n=== Project Management ===\n")

	// Create a project
	project, err := client.Projects().Create(ctx, luna.ProjectCreate{
		Name:        "My Awesome App",
		Description: stringPtr("A revolutionary application built with Luna SDK"),
	})
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	fmt.Printf("Created project: %s\n", project.ID)

	// List all projects
	projects, err := client.Projects().List(ctx, &luna.ListParams{Limit: 20})
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}
	fmt.Printf("Total projects: %d\n", len(projects.Data))

	// Get project details
	projectDetails, err := client.Projects().Get(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}
	fmt.Printf("Project: %s\n", projectDetails.Name)
	fmt.Printf("Owner: %s\n", projectDetails.OwnerID)
	fmt.Printf("Created: %s\n", projectDetails.CreatedAt)

	// Update project
	updated, err := client.Projects().Update(ctx, project.ID, luna.ProjectUpdate{
		Description: stringPtr("Updated description with new features"),
	})
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	fmt.Printf("Updated project: %s\n", *updated.Description)

	// Clean up
	if err := client.Projects().Delete(ctx, project.ID); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	fmt.Println("Project deleted")

	return nil
}

// ============================================
// Example 3: Paginating Through Results
// ============================================

func paginationExample(ctx context.Context, client *luna.Client) error {
	fmt.Println("\n=== Pagination Example ===\n")

	// Using the iterator for automatic pagination
	count := 0
	paginator := client.Users().Iterate(ctx, &luna.ListParams{Limit: 10})

	for paginator.Next() {
		user := paginator.Current()
		fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
		count++
		if count >= 50 { // Limit for demo
			break
		}
	}

	if err := paginator.Err(); err != nil {
		return fmt.Errorf("pagination error: %w", err)
	}

	fmt.Printf("Iterated through %d users\n", count)
	return nil
}

// ============================================
// Example 4: Error Handling
// ============================================

func errorHandlingExample(ctx context.Context, client *luna.Client) error {
	fmt.Println("\n=== Error Handling ===\n")

	// Try to get a non-existent user
	_, err := client.Users().Get(ctx, "usr_nonexistent123")
	if err != nil {
		// Check for specific error types
		switch e := err.(type) {
		case *luna.NotFoundError:
			fmt.Printf("Not found error: %s\n", e.Message)
			fmt.Printf("Error code: %s\n", e.Code)
		case *luna.ValidationError:
			fmt.Printf("Validation error: %s\n", e.Message)
		case *luna.RateLimitError:
			fmt.Printf("Rate limited! Retry after: %d seconds\n", e.RetryAfter)
		case *luna.AuthenticationError:
			fmt.Printf("Authentication error: %s\n", e.Message)
		default:
			fmt.Printf("Unexpected error: %v\n", err)
		}
	}

	return nil
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
