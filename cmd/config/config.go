package config

import (
	"fmt"
	"os"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	ConfigCmd.AddCommand(setCmd)
	ConfigCmd.AddCommand(showCmd)
	setCmd.AddCommand(setTokenCmd)
	setCmd.AddCommand(setEndpointCmd)
}

// ConfigCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure CLI settings (Token, Endpoint, etc.)",
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration value",
}

var setTokenCmd = &cobra.Command{
	Use:   "token <value>",
	Short: "Set the API token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		config.GlobalConfig.APIToken = token
		if err := config.SaveConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("API Token set successfully")
	},
}

var setEndpointCmd = &cobra.Command{
	Use:   "endpoint <value>",
	Short: "Set the API endpoint",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := args[0]
		config.GlobalConfig.APIEndpoint = endpoint
		if err := config.SaveConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("API Endpoint set successfully")
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current Configuration:")
		fmt.Printf("  API Token:    %s\n", config.GlobalConfig.APIToken)
		fmt.Printf("  API Endpoint: %s\n", config.GlobalConfig.APIEndpoint)
		fmt.Printf("  Output Format: %s\n", config.GlobalConfig.OutputFormat)
	},
}
