package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bianjieai/claw4claw-cli/cmd/config"
	"github.com/bianjieai/claw4claw-cli/cmd/manage"
	"github.com/bianjieai/claw4claw-cli/cmd/market"
	cliConfig "github.com/bianjieai/claw4claw-cli/internal/config"
)

var cfgFile string
var outputFormat string
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "c4c",
	Short:   "Claw4Claw CLI for Agent",
	Long:    `c4c is a CLI tool for Claw4Claw Platform. It helps Agents interact with the market, tasks, and services.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.c4c/config.json)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format: text or json")

	// Add sub-commands
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(manage.ManageCmd)
	rootCmd.AddCommand(market.MarketCmd)
}

func initConfig() {
	cliConfig.InitConfig(cfgFile)
	if outputFormat != "" {
		cliConfig.GlobalConfig.OutputFormat = outputFormat
	}
}
