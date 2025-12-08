package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  `Manage projects on the Luna platform.`,
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all projects with pagination support.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")
		_ = limit
		_ = cursor

		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("not authenticated. Run 'luna auth login' or set LUNA_API_KEY")
		}

		// Mock response for demonstration
		projects := []map[string]interface{}{
			{"id": "prj_abc123", "name": "Project Alpha", "owner_id": "usr_xyz", "created_at": "2024-01-15"},
			{"id": "prj_def456", "name": "Project Beta", "owner_id": "usr_xyz", "created_at": "2024-01-16"},
		}

		return outputProjects(projects)
	},
}

var projectsGetCmd = &cobra.Command{
	Use:   "get [project-id]",
	Short: "Get a project by ID",
	Long:  `Retrieve detailed information about a specific project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := args[0]

		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("not authenticated. Run 'luna auth login' or set LUNA_API_KEY")
		}

		// Mock response for demonstration
		project := map[string]interface{}{
			"id":          projectID,
			"name":        "Project Alpha",
			"description": "A sample project",
			"owner_id":    "usr_xyz",
			"created_at":  "2024-01-15",
			"updated_at":  "2024-01-15",
		}

		return outputProject(project)
	},
}

var projectsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long:  `Create a new project with the specified details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("not authenticated. Run 'luna auth login' or set LUNA_API_KEY")
		}

		fmt.Printf("✓ Created project: %s\n", name)
		if description != "" {
			fmt.Printf("  Description: %s\n", description)
		}
		return nil
	},
}

var projectsDeleteCmd = &cobra.Command{
	Use:   "delete [project-id]",
	Short: "Delete a project",
	Long:  `Delete a project by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := args[0]

		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("not authenticated. Run 'luna auth login' or set LUNA_API_KEY")
		}

		fmt.Printf("✓ Deleted project: %s\n", projectID)
		return nil
	},
}

func init() {
	projectsListCmd.Flags().Int("limit", 20, "Maximum number of results")
	projectsListCmd.Flags().String("cursor", "", "Pagination cursor")

	projectsCreateCmd.Flags().String("name", "", "Project name (required)")
	projectsCreateCmd.Flags().String("description", "", "Project description")

	projectsCmd.AddCommand(projectsListCmd)
	projectsCmd.AddCommand(projectsGetCmd)
	projectsCmd.AddCommand(projectsCreateCmd)
	projectsCmd.AddCommand(projectsDeleteCmd)
}

func outputProjects(projects []map[string]interface{}) error {
	switch outputFmt {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(projects)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(projects)
	default: // table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tOWNER\tCREATED")
		for _, p := range projects {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				p["id"], p["name"], p["owner_id"], p["created_at"])
		}
		return w.Flush()
	}
}

func outputProject(project map[string]interface{}) error {
	switch outputFmt {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(project)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(project)
	default: // table
		fmt.Printf("ID:          %s\n", project["id"])
		fmt.Printf("Name:        %s\n", project["name"])
		if desc, ok := project["description"]; ok && desc != "" {
			fmt.Printf("Description: %s\n", desc)
		}
		fmt.Printf("Owner:       %s\n", project["owner_id"])
		fmt.Printf("Created:     %s\n", project["created_at"])
		fmt.Printf("Updated:     %s\n", project["updated_at"])
		return nil
	}
}
