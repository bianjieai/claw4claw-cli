package service

import (
	"fmt"
	"os"

	"github.com/bianjieai/claw4claw-cli/internal/client"
	"github.com/bianjieai/claw4claw-cli/internal/printer"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func HireAgent(req types.CreateEmploymentRequest) error {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.CreateEmployment(req)
	if err != nil {
		return fmt.Errorf("error hiring agent: %w", err)
	}
	printer.PrintCreateEmployment(os.Stdout, resp)
	return nil
}

func FireAgent(employmentID uint, req types.TerminateEmploymentRequest) error {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.TerminateEmployment(employmentID, req)
	if err != nil {
		return fmt.Errorf("error terminating employment: %w", err)
	}
	printer.PrintTerminateEmployment(os.Stdout, resp)
	return nil
}

func GetMyEmployments(params types.MyEmploymentsQueryParams) error {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.GetMyEmployments(params)
	if err != nil {
		return fmt.Errorf("error fetching employments: %w", err)
	}
	printer.PrintEmploymentList(os.Stdout, resp, params.Role)
	return nil
}

func GetEmploymentDetail(employmentID uint) error {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.GetEmploymentByID(employmentID)
	if err != nil {
		return fmt.Errorf("error fetching employment detail: %w", err)
	}
	printer.PrintEmploymentDetail(os.Stdout, resp)
	return nil
}

func AcceptEmployment(employmentID uint, req types.AcceptEmploymentRequest) error {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.AcceptEmployment(employmentID, req)
	if err != nil {
		return fmt.Errorf("error accepting employment: %w", err)
	}
	fmt.Println("Employment accepted successfully")
	fmt.Printf("Employment ID: %d\n", resp.ID)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Started At: %s\n", resp.StartTime)
	return nil
}

func RejectEmployment(employmentID uint, req types.RejectEmploymentRequest) error {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.RejectEmployment(employmentID, req)
	if err != nil {
		return fmt.Errorf("error rejecting employment: %w", err)
	}
	fmt.Println("Employment rejected successfully")
	fmt.Printf("Employment ID: %d\n", resp.ID)
	fmt.Printf("Status: %s\n", resp.Status)
	return nil
}
