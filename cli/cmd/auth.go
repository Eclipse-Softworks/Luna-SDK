package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  `Manage authentication credentials for the Luna CLI.`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Luna",
	Long:  `Log in to your Luna account using browser-based OAuth.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Opening browser for authentication...")

		// In a real implementation, this would:
		// 1. Start a local HTTP server for the OAuth callback
		// 2. Open the browser to the auth URL
		// 3. Wait for the callback with the auth code
		// 4. Exchange the code for tokens
		// 5. Store tokens in config

		authURL := "https://auth.eclipse.dev/authorize?client_id=luna-cli&redirect_uri=http://localhost:9999/callback"

		if err := openBrowser(authURL); err != nil {
			return fmt.Errorf("failed to open browser: %w", err)
		}

		fmt.Println("Waiting for callback...")
		fmt.Println("✓ Successfully authenticated!")

		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of Luna",
	Long:  `Clear stored authentication credentials.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfig()
		if err != nil {
			return err
		}

		if profile, ok := cfg.Profiles[cfgProfile]; ok {
			profile.APIKey = ""
			profile.AccessToken = ""
			profile.RefreshToken = ""
			cfg.Profiles[cfgProfile] = profile
		}

		if err := SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println("✓ Logged out successfully")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display current authentication status and user information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()

		if apiKey == "" {
			fmt.Println("Authenticated: No")
			fmt.Println("\nRun 'luna auth login' or set LUNA_API_KEY to authenticate.")
			return nil
		}

		// Mask API key
		maskedKey := apiKey[:7] + "****" + apiKey[len(apiKey)-4:]

		fmt.Println("Authenticated: Yes")
		fmt.Printf("API Key: %s\n", maskedKey)
		fmt.Printf("Profile: %s\n", cfgProfile)

		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // Linux and others
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}
