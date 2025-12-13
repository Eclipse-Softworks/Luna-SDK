// Package main demonstrates a complete application using Luna SDK.
//
// This example builds a project management CLI tool.
// Run with: go run project_manager.go
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
)

// ProjectManager handles project and team operations
type ProjectManager struct {
	client *luna.Client
}

// NewProjectManager creates a new project manager
func NewProjectManager(client *luna.Client) *ProjectManager {
	return &ProjectManager{client: client}
}

// TeamMember represents a team member
type TeamMember struct {
	UserID string
	Name   string
	Email  string
	Role   string
}

// Team represents a team with its project and members
type Team struct {
	ProjectID string
	Name      string
	Members   []TeamMember
}

func main() {
	client := luna.NewClient(
		luna.WithAPIKey(os.Getenv("LUNA_API_KEY")),
	)

	pm := NewProjectManager(client)
	ctx := context.Background()

	// Run the demo
	if err := runDemo(ctx, pm); err != nil {
		fmt.Printf("Demo failed: %v\n", err)
		os.Exit(1)
	}
}

func runDemo(ctx context.Context, pm *ProjectManager) error {
	// Create a new team/project
	team, err := pm.CreateTeam(ctx, "Project Phoenix", "A revolutionary new product initiative")
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	defer pm.Cleanup(ctx, team) // Clean up at the end

	// Add team members
	members := []struct {
		Email string
		Name  string
		Role  string
	}{
		{"alice@example.com", "Alice Johnson", "Project Lead"},
		{"bob@example.com", "Bob Smith", "Developer"},
		{"carol@example.com", "Carol Williams", "Designer"},
	}

	for _, m := range members {
		if err := pm.AddMember(ctx, team, m.Email, m.Name, m.Role); err != nil {
			fmt.Printf("Warning: Failed to add member %s: %v\n", m.Name, err)
		}
	}

	// Show dashboard
	pm.ShowDashboard(team)

	// Generate AI project summary
	if err := pm.GenerateProjectSummary(ctx, team); err != nil {
		fmt.Printf("Warning: Failed to generate summary: %v\n", err)
	}

	// Get task suggestions
	if err := pm.SuggestTasks(ctx, team, "We're starting a new mobile app. Need to set up dev environment and create designs."); err != nil {
		fmt.Printf("Warning: Failed to suggest tasks: %v\n", err)
	}

	// Show storage info
	if err := pm.ShowStorageInfo(ctx); err != nil {
		fmt.Printf("Warning: Failed to show storage info: %v\n", err)
	}

	// Show available workflows
	if err := pm.ShowWorkflows(ctx); err != nil {
		fmt.Printf("Warning: Failed to show workflows: %v\n", err)
	}

	return nil
}

// ============================================
// Team Management
// ============================================

// CreateTeam creates a new team (project)
func (pm *ProjectManager) CreateTeam(ctx context.Context, name, description string) (*Team, error) {
	fmt.Printf("\nCreating team: %s\n", name)

	project, err := pm.client.Projects().Create(ctx, luna.ProjectCreate{
		Name:        name,
		Description: &description,
	})
	if err != nil {
		return nil, err
	}

	team := &Team{
		ProjectID: project.ID,
		Name:      project.Name,
		Members:   []TeamMember{},
	}

	fmt.Printf("Team created with ID: %s\n", project.ID)
	return team, nil
}

// AddMember adds a member to the team
func (pm *ProjectManager) AddMember(ctx context.Context, team *Team, email, name, role string) error {
	fmt.Printf("\nAdding member: %s\n", name)

	// Create the user
	user, err := pm.client.Users().Create(ctx, luna.UserCreate{
		Email: email,
		Name:  name,
	})
	if err != nil {
		// User might already exist, try to find them
		users, listErr := pm.client.Users().List(ctx, &luna.ListParams{Limit: 100})
		if listErr != nil {
			return fmt.Errorf("failed to create or find user: %w", err)
		}

		for _, u := range users.Data {
			if u.Email == email {
				user = &u
				fmt.Printf("   Found existing user: %s\n", user.ID)
				break
			}
		}

		if user == nil {
			return fmt.Errorf("could not find or create user: %w", err)
		}
	} else {
		fmt.Printf("   Created new user: %s\n", user.ID)
	}

	member := TeamMember{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   role,
	}

	team.Members = append(team.Members, member)
	fmt.Printf("%s added as %s\n", name, role)
	return nil
}

// ListMembers displays all team members
func (pm *ProjectManager) ListMembers(team *Team) {
	fmt.Printf("\nTeam: %s\n", team.Name)
	fmt.Println(strings.Repeat("-", 40))

	if len(team.Members) == 0 {
		fmt.Println("   No members yet")
		return
	}

	for _, member := range team.Members {
		fmt.Printf("   • %s (%s)\n", member.Name, member.Email)
		fmt.Printf("     Role: %s\n", member.Role)
		fmt.Printf("     ID: %s\n", member.UserID)
	}
}

// ============================================
// AI-Powered Features
// ============================================

// GenerateProjectSummary uses AI to generate a project summary
func (pm *ProjectManager) GenerateProjectSummary(ctx context.Context, team *Team) error {
	fmt.Printf("\nGenerating project summary...\n")

	project, err := pm.client.Projects().Get(ctx, team.ProjectID)
	if err != nil {
		return err
	}

	description := "No description"
	if project.Description != nil {
		description = *project.Description
	}

	response, err := pm.client.AI().ChatCompletions(ctx, luna.CompletionRequest{
		Model: "luna-gpt-4",
		Messages: []luna.Message{
			{
				Role:    "system",
				Content: "You are a project manager assistant. Generate concise, professional summaries.",
			},
			{
				Role: "user",
				Content: fmt.Sprintf(`Generate a brief project status summary for:
				
Project: %s
Description: %s
Team Size: %d members
Created: %s

Include:
1. Project overview
2. Team composition
3. Suggested next steps`, project.Name, description, len(team.Members), project.CreatedAt),
			},
		},
		Temperature: 0.5,
	})
	if err != nil {
		return err
	}

	fmt.Println("\nProject Summary:\n")
	fmt.Println(response.Choices[0].Message.Content)
	return nil
}

// SuggestTasks uses AI to suggest tasks for the team
func (pm *ProjectManager) SuggestTasks(ctx context.Context, team *Team, context string) error {
	fmt.Printf("\nGenerating task suggestions...\n")

	memberNames := make([]string, len(team.Members))
	for i, m := range team.Members {
		memberNames[i] = m.Name
	}

	response, err := pm.client.AI().ChatCompletions(ctx, luna.CompletionRequest{
		Model: "luna-gpt-4",
		Messages: []luna.Message{
			{
				Role:    "system",
				Content: "You are a project planning assistant. Suggest actionable tasks.",
			},
			{
				Role: "user",
				Content: fmt.Sprintf(`Based on this context, suggest 5 actionable tasks:
				
Team: %s
Team Members: %s
Context: %s

Format as a numbered list with assigned team member (if applicable).`,
					team.Name, strings.Join(memberNames, ", "), context),
			},
		},
		Temperature: 0.7,
	})
	if err != nil {
		return err
	}

	fmt.Println("\nSuggested Tasks:\n")
	fmt.Println(response.Choices[0].Message.Content)
	return nil
}

// ============================================
// Storage Management
// ============================================

// ShowStorageInfo displays storage bucket information
func (pm *ProjectManager) ShowStorageInfo(ctx context.Context) error {
	fmt.Printf("\nStorage Information\n")
	fmt.Println(strings.Repeat("-", 40))

	buckets, err := pm.client.Storage().Buckets().List(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Available buckets: %d\n", len(buckets.Data))

	for _, bucket := range buckets.Data {
		fmt.Printf("\n   [BUCKET] %s\n", bucket.Name)
		fmt.Printf("      ID: %s\n", bucket.ID)
		fmt.Printf("      Region: %s\n", bucket.Region)

		// List files in bucket
		files, err := pm.client.Storage().Files().List(ctx, bucket.ID)
		if err != nil {
			fmt.Printf("      Files: (error listing)\n")
			continue
		}
		fmt.Printf("      Files: %d\n", len(files.Data))
	}

	return nil
}

// ============================================
// Workflow Management
// ============================================

// ShowWorkflows displays available workflows
func (pm *ProjectManager) ShowWorkflows(ctx context.Context) error {
	fmt.Printf("\nAvailable Workflows\n")
	fmt.Println(strings.Repeat("-", 40))

	workflows, err := pm.client.Automation().Workflows().List(ctx)
	if err != nil {
		return err
	}

	activeCount := 0
	for _, w := range workflows.Data {
		if w.IsActive {
			activeCount++
		}
	}

	fmt.Printf("Total workflows: %d (%d active)\n\n", len(workflows.Data), activeCount)

	for _, workflow := range workflows.Data {
		status := "[Inactive]"
		if workflow.IsActive {
			status = "[Active]"
		}

		fmt.Printf("   %s %s\n", status, workflow.Name)
		fmt.Printf("      Trigger: %s\n", workflow.TriggerType)
		fmt.Printf("      ID: %s\n", workflow.ID)
	}

	return nil
}

// ============================================
// Dashboard
// ============================================

// ShowDashboard displays a team dashboard
func (pm *ProjectManager) ShowDashboard(team *Team) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("TEAM DASHBOARD: %s\n", team.Name)
	fmt.Println(strings.Repeat("=", 60))

	// Team info
	fmt.Println("\nOverview")
	fmt.Printf("   Project ID: %s\n", team.ProjectID)
	fmt.Printf("   Members: %d\n", len(team.Members))

	// Members
	pm.ListMembers(team)

	// Recent activity (simulated)
	fmt.Println("\nRecent Activity")
	fmt.Printf("   • Team created at %s\n", time.Now().Format("15:04:05"))
	for _, member := range team.Members {
		fmt.Printf("   - %s joined the team\n", member.Name)
	}

	// Quick actions
	fmt.Println("\nQuick Actions")
	fmt.Println("   1. Add new member")
	fmt.Println("   2. Generate project summary")
	fmt.Println("   3. Suggest tasks")
	fmt.Println("   4. Upload file")
	fmt.Println("   5. Trigger workflow")
}

// ============================================
// Cleanup
// ============================================

// Cleanup removes all created resources
func (pm *ProjectManager) Cleanup(ctx context.Context, team *Team) {
	fmt.Println("\nCleaning up...")

	for _, member := range team.Members {
		if err := pm.client.Users().Delete(ctx, member.UserID); err != nil {
			fmt.Printf("   Warning: Failed to delete user %s: %v\n", member.Name, err)
		} else {
			fmt.Printf("   Deleted user: %s\n", member.Name)
		}
	}

	if err := pm.client.Projects().Delete(ctx, team.ProjectID); err != nil {
		fmt.Printf("   Warning: Failed to delete project: %v\n", err)
	} else {
		fmt.Printf("   Deleted project: %s\n", team.Name)
	}

	fmt.Println("Cleanup complete")
}
