package manage

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/bianjieai/claw4claw-cli/internal/service"
	"github.com/bianjieai/claw4claw-cli/internal/types"
	"github.com/spf13/cobra"
)

var (
	invocationRole   string
	invocationStatus string
	invocationPage   int
	invocationLimit  int

	invocationInputFile      string
	invocationMaxPrice       float64
	invocationIdempotencyKey string

	invocationSubmitStatus      string
	invocationSubmitOutputFile  string
	invocationSubmitError       string
	invocationSubmitAttachments []string

	invocationReviewRating int
	invocationReviewText   string
)

func init() {
	ManageCmd.AddCommand(invocationCmd)
	invocationCmd.AddCommand(invocationListCmd)
	invocationCmd.AddCommand(invocationShowCmd)
	invocationCmd.AddCommand(invocationInvokeCmd)
	invocationCmd.AddCommand(invocationSubmitCmd)
	invocationCmd.AddCommand(invocationReviewCmd)

	invocationListCmd.Flags().StringVarP(&invocationRole, "role", "r", "", "Role filter (caller/provider)")
	invocationListCmd.Flags().StringVarP(&invocationStatus, "status", "s", "", "Status filter (pending/completed/failed/timeout)")
	invocationListCmd.Flags().IntVarP(&invocationPage, "page", "p", 1, "Page number")
	invocationListCmd.Flags().IntVarP(&invocationLimit, "limit", "l", 10, "Items per page")

	invocationInvokeCmd.Flags().StringVarP(&invocationInputFile, "input", "i", "", "JSON file containing input parameters")
	invocationInvokeCmd.Flags().Float64Var(&invocationMaxPrice, "max-price", 0, "Maximum price")
	invocationInvokeCmd.Flags().StringVar(&invocationIdempotencyKey, "idempotency-key", "", "Idempotency key for deduplication")

	invocationSubmitCmd.Flags().StringVarP(&invocationSubmitStatus, "status", "s", "completed", "Result status (completed/failed)")
	invocationSubmitCmd.Flags().StringVarP(&invocationSubmitOutputFile, "output", "o", "", "JSON file containing output")
	invocationSubmitCmd.Flags().StringVarP(&invocationSubmitError, "error", "e", "", "Error message (if status is failed)")
	invocationSubmitCmd.Flags().StringArrayVarP(&invocationSubmitAttachments, "attachment", "a", nil, "Attachment URLs (can be specified multiple times)")

	invocationReviewCmd.Flags().IntVarP(&invocationReviewRating, "rating", "r", 0, "Rating (1-5)")
	invocationReviewCmd.Flags().StringVarP(&invocationReviewText, "review", "t", "", "Review text")
}

var invocationCmd = &cobra.Command{
	Use:     "service-invocation",
	Aliases: []string{"invocation", "inv"},
	Short:   "Manage service invocations",
}

var invocationListCmd = &cobra.Command{
	Use:   "list",
	Short: "Caller/Provider: List your service invocations (use --role caller/provider)",
	Run: func(cmd *cobra.Command, args []string) {
		params := service.GetServiceInvocationsParams{
			Role:   invocationRole,
			Status: invocationStatus,
			Page:   invocationPage,
			Limit:  invocationLimit,
		}
		service.GetServiceInvocationList(params)
	},
}

var invocationShowCmd = &cobra.Command{
	Use:   "show <invocation-id>",
	Short: "Caller/Provider: Show details of a service invocation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid invocation ID: %v\n", err)
			os.Exit(1)
		}
		service.GetServiceInvocationDetail(id)
	},
}

var invocationInvokeCmd = &cobra.Command{
	Use:   "invoke <service-id>",
	Short: "Caller: Invoke a service from market",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid service ID: %v\n", err)
			os.Exit(1)
		}

		var input map[string]interface{}

		if invocationInputFile != "" {
			data, err := os.ReadFile(invocationInputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
				os.Exit(1)
			}
			if err := json.Unmarshal(data, &input); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing input JSON: %v\n", err)
				os.Exit(1)
			}
		}

		req := types.InvokeServiceRequest{
			ServiceID:      serviceID,
			Input:          input,
			IdempotencyKey: invocationIdempotencyKey,
		}

		if invocationMaxPrice > 0 {
			req.MaxPrice = &invocationMaxPrice
		}

		service.InvokeService(req)
	},
}

var invocationSubmitCmd = &cobra.Command{
	Use:   "submit <invocation-id>",
	Short: "Provider: Submit result for a service invocation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid invocation ID: %v\n", err)
			os.Exit(1)
		}

		var output map[string]interface{}

		if invocationSubmitOutputFile != "" {
			data, err := os.ReadFile(invocationSubmitOutputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading output file: %v\n", err)
				os.Exit(1)
			}
			if err := json.Unmarshal(data, &output); err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing output JSON: %v\n", err)
				os.Exit(1)
			}
		}

		req := types.SubmitServiceResultRequest{
			Status:      invocationSubmitStatus,
			Output:      output,
			Error:       invocationSubmitError,
			Attachments: invocationSubmitAttachments,
		}

		service.SubmitServiceResult(id, req)
	},
}

var invocationReviewCmd = &cobra.Command{
	Use:   "review <invocation-id>",
	Short: "Caller: Review a service invocation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid invocation ID: %v\n", err)
			os.Exit(1)
		}

		req := types.ReviewServiceInvocationRequest{
			Review: invocationReviewText,
		}

		if invocationReviewRating > 0 {
			req.Rating = &invocationReviewRating
		}

		service.ReviewServiceInvocation(id, req)
	},
}
