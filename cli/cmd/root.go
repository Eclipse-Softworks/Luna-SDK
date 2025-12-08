// Package cmd provides CLI commands for the Luna SDK.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	cfgProfile string
	apiKey     string
	outputFmt  string
	noColor    bool
	verbose    bool
	debug      bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "luna",
	Short: "Luna CLI - Eclipse Softworks Platform",
	Long: `Luna CLI provides a command-line interface to the Eclipse Softworks Platform.

Use Luna CLI to manage users, projects, and other platform resources
directly from your terminal.

Examples:
  luna auth login
  luna users list
  luna projects create --name "My Project"`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgProfile, "profile", "default", "Configuration profile to use")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "format", "f", "table", "Output format: table, json, yaml")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")

	// Add subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(configCmd)
}

// getAPIKey returns the API key from flag or config
func getAPIKey() string {
	if apiKey != "" {
		return apiKey
	}

	// Try environment variable
	if envKey := os.Getenv("LUNA_API_KEY"); envKey != "" {
		return envKey
	}

	// Try config file
	cfg, err := LoadConfig()
	if err == nil {
		if profile, ok := cfg.Profiles[cfgProfile]; ok {
			return profile.APIKey
		}
	}

	return ""
}

// printError prints an error message
func printError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}
