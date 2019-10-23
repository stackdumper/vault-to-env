package command

import (
	log "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
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
		log.WithField("leases", revokeFlags.leases).Debug("revoking leases")

		for _, lease := range revokeFlags.leases {
			logger := log.WithField("lease", lease)
			logger.Debug("revoking lease")

			err := PersistentState.client.RevokeLease(lease)
			if err != nil {
				logger.WithError(err).Fatal("failed to revoke lease")
			}
		}
	},
}
