package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Config represents the CLI configuration file
type Config struct {
	DefaultProfile string             `yaml:"default_profile"`
	Profiles       map[string]Profile `yaml:"profiles"`
	Settings       Settings           `yaml:"settings"`
}

// Profile represents a configuration profile
type Profile struct {
	APIKey       string `yaml:"api_key,omitempty"`
	BaseURL      string `yaml:"base_url,omitempty"`
	AccessToken  string `yaml:"access_token,omitempty"`
	RefreshToken string `yaml:"refresh_token,omitempty"`
}

// Settings holds CLI settings
type Settings struct {
	OutputFormat string `yaml:"output_format"`
	Color        bool   `yaml:"color"`
	Pager        string `yaml:"pager"`
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `Manage configuration settings for the Luna CLI.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Examples:
  luna config set api_key lk_live_xxxx
  luna config set base_url https://api.staging.eclipse.dev`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := LoadConfig()
		if err != nil {
			cfg = &Config{
				DefaultProfile: "default",
				Profiles:       make(map[string]Profile),
				Settings: Settings{
					OutputFormat: "table",
					Color:        true,
				},
			}
		}

		if _, ok := cfg.Profiles[cfgProfile]; !ok {
			cfg.Profiles[cfgProfile] = Profile{}
		}

		profile := cfg.Profiles[cfgProfile]

		switch key {
		case "api_key":
			profile.APIKey = value
		case "base_url":
			profile.BaseURL = value
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		cfg.Profiles[cfgProfile] = profile

		if err := SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✓ Set %s for profile '%s'\n", key, cfgProfile)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  `Get a configuration value.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		cfg, err := LoadConfig()
		if err != nil {
			return fmt.Errorf("no configuration found")
		}

		profile, ok := cfg.Profiles[cfgProfile]
		if !ok {
			return fmt.Errorf("profile '%s' not found", cfgProfile)
		}

		var value string
		switch key {
		case "api_key":
			if profile.APIKey != "" {
				value = profile.APIKey[:7] + "****" + profile.APIKey[len(profile.APIKey)-4:]
			}
		case "base_url":
			value = profile.BaseURL
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		if value == "" {
			fmt.Printf("%s: (not set)\n", key)
		} else {
			fmt.Printf("%s: %s\n", key, value)
		}

		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration",
	Long:  `List all configuration settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfig()
		if err != nil {
			fmt.Println("No configuration found.")
			fmt.Println("Run 'luna config set api_key <key>' to configure.")
			return nil
		}

		fmt.Printf("Default Profile: %s\n\n", cfg.DefaultProfile)

		for name, profile := range cfg.Profiles {
			fmt.Printf("[%s]\n", name)
			if profile.APIKey != "" {
				masked := profile.APIKey[:7] + "****" + profile.APIKey[len(profile.APIKey)-4:]
				fmt.Printf("  api_key: %s\n", masked)
			}
			if profile.BaseURL != "" {
				fmt.Printf("  base_url: %s\n", profile.BaseURL)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".luna", "config.yaml")
}

// LoadConfig loads the configuration from disk
func LoadConfig() (*Config, error) {
	path := getConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves the configuration to disk
func SaveConfig(cfg *Config) error {
	path := getConfigPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
