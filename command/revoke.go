package command

import (
	"log"

	"github.com/spf13/cobra"
)

var renewFlags struct {
	leases   []string
	duration int
}

func init() {
	renewCmd.Flags().StringArrayVar(&renewFlags.leases, "leases", []string{}, "list leases to renew")
	renewCmd.Flags().IntVar(&renewFlags.duration, "duration", 3600, "lease renew duration")
}

var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "Renew secrets leases",
	Run: func(cmd *cobra.Command, args []string) {
		for _, lease := range renewFlags.leases {
			err := PersistentState.client.RenewLease(lease, renewFlags.duration)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
