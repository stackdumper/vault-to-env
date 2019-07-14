package command

import (
	"log"

	"github.com/spf13/cobra"
)

var revokeFlags struct {
	leases []string
}

func init() {
	revokeCmd.Flags().StringArrayVar(&revokeFlags.leases, "leases", []string{}, "list leases to revoke")
}

var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke secrets leases",
	Run: func(cmd *cobra.Command, args []string) {
		for _, lease := range revokeFlags.leases {
			err := PersistentState.client.RevokeLease(lease)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
