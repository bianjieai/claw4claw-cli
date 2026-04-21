package printer

import (
	"fmt"
	"io"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func PrintWalletSummary(w io.Writer, summary *types.WalletSummary) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, summary)
		return
	}

	fmt.Fprintln(w, "Wallet Summary")
	fmt.Fprintln(w, "==============")
	fmt.Fprintf(w, "Available Balance:  %d Shells (¥%.2f)\n", summary.Balance, summary.BalanceInCNY)
	fmt.Fprintf(w, "Frozen Amount:      %d Shells\n", summary.FrozenAmount)
	fmt.Fprintf(w, "Total Assets:       %d Shells\n", summary.TotalAssets)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Invoicable Amount:  ¥%.2f\n", summary.InvoicableAmount)
	fmt.Fprintf(w, "Pending Refunds:    %d\n", summary.PendingRefunds)
}
