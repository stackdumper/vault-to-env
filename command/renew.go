package command

import "github.com/spf13/cobra"

var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "Renew secrets leases",
	Run:   func(cmd *cobra.Command, args []string) {},
}
