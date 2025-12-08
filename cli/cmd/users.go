package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
	Long:  `Manage users on the Luna platform.`,
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `List all users with pagination support.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")
		_ = limit
		_ = cursor

		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("not authenticated. Run 'luna auth login' or set LUNA_API_KEY")
		}

		client, err := luna.NewClient(luna.WithAPIKey(apiKey))
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		users, err := client.Users().List(cmd.Context(), &luna.ListParams{
			Limit:  limit,
			Cursor: cursor,
		})
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		// Convert to map for output compatibility (or update output function)
		// For now we map strictly to the output format expected
		var output []map[string]interface{}
		// Re-marshal to map for generic output handling
		data, _ := json.Marshal(users.Data)
		json.Unmarshal(data, &output)

		return outputUsers(output)
	},
}

var usersGetCmd = &cobra.Command{
	Use:   "get [user-id]",
	Short: "Get a user by ID",
	Long:  `Retrieve detailed information about a specific user.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]

		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("not authenticated. Run 'luna auth login' or set LUNA_API_KEY")
		}

		client, err := luna.NewClient(luna.WithAPIKey(apiKey))
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		user, err := client.Users().Get(cmd.Context(), userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		var output map[string]interface{}
		data, _ := json.Marshal(user)
		json.Unmarshal(data, &output)

		return outputUser(output)
	},
}

var usersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long:  `Create a new user with the specified details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")

		if name == "" || email == "" {
			return fmt.Errorf("--name and --email are required")
		}

		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("not authenticated. Run 'luna auth login' or set LUNA_API_KEY")
		}

		client, err := luna.NewClient(luna.WithAPIKey(apiKey))
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		user, err := client.Users().Create(cmd.Context(), luna.UserCreate{
			Name:  name,
			Email: email,
		})
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		fmt.Printf("✓ Created user: %s <%s> (%s)\n", user.Name, user.Email, user.ID)
		return nil
	},
}

var usersDeleteCmd = &cobra.Command{
	Use:   "delete [user-id]",
	Short: "Delete a user",
	Long:  `Delete a user by their ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]

		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("not authenticated. Run 'luna auth login' or set LUNA_API_KEY")
		}

		client, err := luna.NewClient(luna.WithAPIKey(apiKey))
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		err = client.Users().Delete(cmd.Context(), userID)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		fmt.Printf("✓ Deleted user: %s\n", userID)
		return nil
	},
}

func init() {
	usersListCmd.Flags().Int("limit", 20, "Maximum number of results")
	usersListCmd.Flags().String("cursor", "", "Pagination cursor")

	usersCreateCmd.Flags().String("name", "", "User name (required)")
	usersCreateCmd.Flags().String("email", "", "User email (required)")

	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersGetCmd)
	usersCmd.AddCommand(usersCreateCmd)
	usersCmd.AddCommand(usersDeleteCmd)
}

func outputUsers(users []map[string]interface{}) error {
	switch outputFmt {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(users)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(users)
	default: // table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tEMAIL\tCREATED")
		for _, u := range users {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				u["id"], u["name"], u["email"], u["created_at"])
		}
		return w.Flush()
	}
}

func outputUser(user map[string]interface{}) error {
	switch outputFmt {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(user)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(user)
	default: // table
		fmt.Printf("ID:         %s\n", user["id"])
		fmt.Printf("Name:       %s\n", user["name"])
		fmt.Printf("Email:      %s\n", user["email"])
		fmt.Printf("Created:    %s\n", user["created_at"])
		fmt.Printf("Updated:    %s\n", user["updated_at"])
		return nil
	}
}
