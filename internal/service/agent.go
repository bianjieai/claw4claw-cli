package service

import (
	"fmt"
	"os"

	"github.com/bianjieai/claw4claw-cli/internal/client"
	"github.com/bianjieai/claw4claw-cli/internal/printer"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func RegisterAgent(req types.RegisterAgentReq) {
	apiClient := client.NewAPIClient()
	agent, err := apiClient.RegisterAgent(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error registering agent: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Agent registered successfully: %s (ID: %d)\n", agent.Name, agent.ID)
}

func GetMyInfo() {
	apiClient := client.NewAPIClient()
	agent, err := apiClient.GetMyInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching agent info: %v\n", err)
		os.Exit(1)
	}
	printer.PrintAgentInfo(os.Stdout, agent)
}

func UpdateMyInfo(req types.UpdateAgentReq) {
	apiClient := client.NewAPIClient()
	err := apiClient.UpdateMyInfo(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating agent: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Agent updated successfully")
}

func SetMyStatus(status string) {
	apiClient := client.NewAPIClient()
	err := apiClient.SetMyStatus(status)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting agent status: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Agent status set to: %s\n", status)
}

func GetMarketAgentList() {
	apiClient := client.NewAPIClient()
	list, err := apiClient.GetMarketAgents()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching market agents: %v\n", err)
		os.Exit(1)
	}
	printer.PrintMarketAgentList(os.Stdout, list)
}

func GetMarketAgentDetail(id string) {
	apiClient := client.NewAPIClient()
	agent, err := apiClient.GetMarketAgentDetail(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching market agent detail: %v\n", err)
		os.Exit(1)
	}
	printer.PrintMarketAgentDetail(os.Stdout, agent)
}

func PublishAgent(req types.PublishAgentReq) {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.PublishAgent(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error publishing agent: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Agent published to market successfully")
	fmt.Printf("Market Visibility: %s\n", resp.MarketVisibility)
	fmt.Printf("Published At: %s\n", resp.MarketPublishedAt)
}

func UnpublishAgent() {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.UnpublishAgent()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unpublishing agent: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Agent unpublished from market successfully")
	fmt.Printf("Market Visibility: %s\n", resp.MarketVisibility)
}

func GetMyBudget() {
	apiClient := client.NewAPIClient()
	budget, err := apiClient.GetMyBudget()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching agent budget: %v\n", err)
		os.Exit(1)
	}
	printer.PrintBudgetInfo(os.Stdout, budget)
}
