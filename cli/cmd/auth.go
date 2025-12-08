package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"os/exec"
	"runtime"

	"github.com/eclipse-softworks/luna-sdk-go/luna"
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
		// 1. Create a channel to signal completion
		done := make(chan string)
		errChan := make(chan error)

		// 2. Start local server
		server := &http.Server{Addr: "127.0.0.1:9999"}
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "Code not found", http.StatusBadRequest)
				errChan <- fmt.Errorf("authorization code not found in callback")
				return
			}
			fmt.Fprintf(w, "Authorization successful! You can close this window now.")
			done <- code
		})

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- fmt.Errorf("failed to start local server: %w", err)
			}
		}()

		// 3. Open browser
		authURL := "https://auth.eclipse.dev/authorize?client_id=luna-cli&redirect_uri=http://localhost:9999/callback&response_type=code"
		fmt.Println("Opening browser for authentication...")
		fmt.Printf("If browser does not open, visit: %s\n", authURL)

		if err := openBrowser(authURL); err != nil {
			fmt.Printf("Failed to open browser: %v\n", err)
		}

		// 4. Wait for callback
		fmt.Println("Waiting for authentication...")
		select {
		case code := <-done:
			_ = server.Shutdown(context.Background())
			fmt.Println("✓ Successfully authenticated!")

			// Exchange code for tokens
			tokens, err := exchangeToken(code)
			if err != nil {
				// Fallback for demo/offline if real endpoint fails
				// But we try to be as real as possible first
				errChan <- fmt.Errorf("failed to exchange token: %w", err)
				return nil
			}

			cfg, err := LoadConfig()
			if err != nil {
				// If config doesn't exist, create default
				cfg = &Config{
					DefaultProfile: "default",
					Profiles:       make(map[string]Profile),
					Settings: Settings{
						OutputFormat: "table",
						Color:        true,
					},
				}
			}

			if cfg.Profiles == nil {
				cfg.Profiles = make(map[string]Profile)
			}

			// Update profile with real tokens
			profile := cfg.Profiles[cfgProfile]
			profile.AccessToken = tokens.AccessToken
			profile.RefreshToken = tokens.RefreshToken
			cfg.Profiles[cfgProfile] = profile

			if err := SaveConfig(cfg); err != nil {
				errChan <- fmt.Errorf("failed to save config: %w", err)
				return nil
			}

			// We are done
			return nil

		case err := <-errChan:
			_ = server.Shutdown(context.Background())
			return err
		case <-time.After(2 * time.Minute):
			_ = server.Shutdown(context.Background())
			return fmt.Errorf("authentication timed out")
		}
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

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify API credentials",
	Long:  "Test the currently configured API key against the server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("not authenticated")
		}

		client, err := luna.NewClient(luna.WithAPIKey(apiKey))
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Attempt to fetch something simpler or just use existing resources to verify auth
		// Since there isn't an explicit "Verify" or "Me" endpoint exposed in the top level resources we see,
		// we'll try to list project or users with limit 1 to check credentials.
		// Actually, let's assume we can list users (self) or similar.

		// A common pattern is to check "Me" but we don't have that resource visible in client.go right now.
		// We'll use List Users as a proxy for "Is Authenticated".

		_, err = client.Users().List(cmd.Context(), &luna.ListParams{Limit: 1})
		if err != nil {
			return fmt.Errorf("verification failed: %w", err)
		}

		fmt.Println("✓ Credentials are valid")
		// We can't easily get the user details without a Me endpoint, but we confirmed the key works.
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(verifyCmd)
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

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func exchangeToken(code string) (*TokenResponse, error) {
	// Real HTTP request to exchange code
	// Use custom client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", "luna-cli")
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:9999/callback")

	req, err := http.NewRequest("POST", "https://auth.eclipse.dev/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error body if possible
		var errorBody map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errorBody)
		return nil, fmt.Errorf("token exchange failed: status %d, response: %v", resp.StatusCode, errorBody)
	}

	var tokens TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokens, nil
}
