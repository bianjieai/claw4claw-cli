package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	OutputFormat string `json:"output"` // "json" or "text"
	APIToken     string `json:"api_token"`
	APIEndpoint  string `json:"api_endpoint"`
	WebhookURL   string `json:"webhook_url"` // Local webhook URL for forwarding messages
}

var GlobalConfig Config
var ConfigFilePath string

// InitConfig initializes configuration from file and environment variables.
func InitConfig(cfgFile string) {
	if cfgFile != "" {
		ConfigFilePath = cfgFile
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configDir := filepath.Join(home, ".c4c")
		ConfigFilePath = filepath.Join(configDir, "config.json")
	}

	// Default configurations
	GlobalConfig = Config{
		OutputFormat: "json",
		APIEndpoint:  "https://api.claw4claw.bianjie.ai",
	}

	// Read from config file
	data, err := os.ReadFile(ConfigFilePath)
	if err == nil {
		if err := json.Unmarshal(data, &GlobalConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling config: %v\n", err)
		}
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
	}

	// Override with environment variables (C4C_ prefix)
	if envToken := os.Getenv("C4C_API_TOKEN"); envToken != "" {
		GlobalConfig.APIToken = envToken
	}
	if envEndpoint := os.Getenv("C4C_API_ENDPOINT"); envEndpoint != "" {
		GlobalConfig.APIEndpoint = envEndpoint
	}
	if envOutput := os.Getenv("C4C_OUTPUT"); envOutput != "" {
		GlobalConfig.OutputFormat = envOutput
	}
	if envWebhook := os.Getenv("C4C_WEBHOOK_URL"); envWebhook != "" {
		GlobalConfig.WebhookURL = envWebhook
	}
}

// SaveConfig saves the current configuration to the config file.
func SaveConfig() error {
	configDir := filepath.Dir(ConfigFilePath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return fmt.Errorf("error creating config directory: %v", err)
		}
	}

	data, err := json.MarshalIndent(&GlobalConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling config: %v", err)
	}

	if err := os.WriteFile(ConfigFilePath, data, 0600); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}
