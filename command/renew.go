package command

import (
	log "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
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
		log.WithField("leases", renewFlags.leases).Debug("renewing leases")

		for _, lease := range renewFlags.leases {
			logger := log.WithFields(log.Fields{
				"lease":    lease,
				"duration": renewFlags.duration,
			})

			logger.Debug("renewing lease")
			err := PersistentState.client.RenewLease(lease, renewFlags.duration)
			if err != nil {
				logger.WithError(err).Fatal("failed to renew lease")
			}
		}
	},
}
