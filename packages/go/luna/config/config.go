// Package config provides configuration loading for the Luna SDK.
package config

import (
	"os"
	"strconv"
)

// Config holds Luna SDK configuration.
type Config struct {
	APIKey       string
	AccessToken  string
	RefreshToken string
	BaseURL      string
	Timeout      int
	MaxRetries   int
	LogLevel     string
}

// EnvVars defines environment variable names.
var EnvVars = struct {
	APIKey       string
	AccessToken  string
	RefreshToken string
	BaseURL      string
	Timeout      string
	MaxRetries   string
	LogLevel     string
}{
	APIKey:       "LUNA_API_KEY",
	AccessToken:  "LUNA_ACCESS_TOKEN",
	RefreshToken: "LUNA_REFRESH_TOKEN",
	BaseURL:      "LUNA_BASE_URL",
	Timeout:      "LUNA_TIMEOUT",
	MaxRetries:   "LUNA_MAX_RETRIES",
	LogLevel:     "LUNA_LOG_LEVEL",
}

// Defaults provides default configuration values.
var Defaults = Config{
	BaseURL:    "https://api.eclipse.dev",
	Timeout:    30000,
	MaxRetries: 3,
	LogLevel:   "info",
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv() Config {
	config := Config{}

	if apiKey := os.Getenv(EnvVars.APIKey); apiKey != "" {
		config.APIKey = apiKey
	}

	if accessToken := os.Getenv(EnvVars.AccessToken); accessToken != "" {
		config.AccessToken = accessToken
	}

	if refreshToken := os.Getenv(EnvVars.RefreshToken); refreshToken != "" {
		config.RefreshToken = refreshToken
	}

	if baseURL := os.Getenv(EnvVars.BaseURL); baseURL != "" {
		config.BaseURL = baseURL
	}

	if timeoutStr := os.Getenv(EnvVars.Timeout); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil && timeout > 0 {
			config.Timeout = timeout
		}
	}

	if maxRetriesStr := os.Getenv(EnvVars.MaxRetries); maxRetriesStr != "" {
		if maxRetries, err := strconv.Atoi(maxRetriesStr); err == nil && maxRetries >= 0 {
			config.MaxRetries = maxRetries
		}
	}

	if logLevel := os.Getenv(EnvVars.LogLevel); logLevel != "" {
		config.LogLevel = logLevel
	}

	return config
}

// Merge merges user config with environment and defaults.
func Merge(userConfig Config) Config {
	envConfig := LoadFromEnv()

	result := Defaults

	// Apply env config
	if envConfig.APIKey != "" {
		result.APIKey = envConfig.APIKey
	}
	if envConfig.AccessToken != "" {
		result.AccessToken = envConfig.AccessToken
	}
	if envConfig.RefreshToken != "" {
		result.RefreshToken = envConfig.RefreshToken
	}
	if envConfig.BaseURL != "" {
		result.BaseURL = envConfig.BaseURL
	}
	if envConfig.Timeout != 0 {
		result.Timeout = envConfig.Timeout
	}
	if envConfig.MaxRetries != 0 {
		result.MaxRetries = envConfig.MaxRetries
	}
	if envConfig.LogLevel != "" {
		result.LogLevel = envConfig.LogLevel
	}

	// Apply user config (takes precedence)
	if userConfig.APIKey != "" {
		result.APIKey = userConfig.APIKey
	}
	if userConfig.AccessToken != "" {
		result.AccessToken = userConfig.AccessToken
	}
	if userConfig.RefreshToken != "" {
		result.RefreshToken = userConfig.RefreshToken
	}
	if userConfig.BaseURL != "" {
		result.BaseURL = userConfig.BaseURL
	}
	if userConfig.Timeout != 0 {
		result.Timeout = userConfig.Timeout
	}
	if userConfig.MaxRetries != 0 {
		result.MaxRetries = userConfig.MaxRetries
	}
	if userConfig.LogLevel != "" {
		result.LogLevel = userConfig.LogLevel
	}

	return result
}
