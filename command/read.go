package command

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"
)

var saveFlags struct {
	vars          []string
	leaseDuration int
	saveLeases    bool
}

func init() {
	saveCmd.Flags().StringSliceVar(&saveFlags.vars, "vars", []string{}, "list of vars to read")
	saveCmd.Flags().IntVar(&saveFlags.leaseDuration, "lease-duration", 0, "adjust secret lease duration")
	saveCmd.Flags().BoolVar(&saveFlags.saveLeases, "save-leases", false, "save secret leases")
}

type ResultVar struct {
	// Input: "Name=Path#Key"
	// Output: "export Name=Value"
	Name    string
	Path    string
	Value   string
	LeaseID string
	Key     []string
}

// Read secrets and output them as env variables
var saveCmd = &cobra.Command{
	Use:   "read",
	Short: "Read secrets and output them as env variables",
	Run: func(cmd *cobra.Command, args []string) {
		var resultVars = make([]ResultVar, len(saveFlags.vars))

		log.WithField("vars", saveFlags.vars).Debug("parsing output variables")
		for index, rawVar := range saveFlags.vars {
			// Var=Path#Key
			iv, ip, ik := getByteIndex(rawVar, '='), getByteIndex(rawVar, '#'), len(rawVar)

			// exit if provided var is incompliete
			if len(rawVar) <= ip+1 {
				log.WithField("variable", rawVar).Fatal("unrecognized variable format, required var=path#key")
			}

			resultVars[index] = ResultVar{
				Name: rawVar[:iv],
				Path: rawVar[iv+1 : ip],
				Key:  strings.Split(rawVar[ip+1:ik], "."),
			}

			log.WithField("var", resultVars[index]).Debug("parsed output variable")
		}

		// get Values
		log.WithField("secrets", resultVars).Debug("reading secrets")
		for i, v := range resultVars {
			result := PersistentState.client.Read(v.Path, v.Key)
			logger := log.WithFields(log.Fields{
				"name": v.Name,
				"path": v.Path,
				"key":  v.Key,
			})

			// print warnings if there are some
			if len(result.Warnings) != 0 {
				for _, warning := range result.Warnings {
					logger.WithField("warning", warning).Warn("received warning while reading secret")
				}
			}

			// exit with error if there's one
			if result.Error != nil {
				logger.WithError(result.Error).Fatal("received error while reading secrets")
			}

			if saveFlags.leaseDuration != 0 && result.LeaseID != "" {
				// adjust lease
				// https://www.vaultproject.io/docs/concepts/lease.html
				err := PersistentState.client.RenewLease(result.LeaseID, saveFlags.leaseDuration)
				if err != nil {
					logger.WithError(result.Error).Fatal("received error while adjusting secret lease")
				}
			}

			// assign value
			resultVars[i].Value = result.Value
			resultVars[i].LeaseID = result.LeaseID

			log.WithField("secret", resultVars[i]).Debug("successfully read secret")
		}

		// extract result string
		// export Var="Value"
		log.Debug("writing result")
		var result string
		for _, v := range resultVars {
			result += fmt.Sprintf("export %s=\"%s\"\n", v.Name, v.Value)

			if saveFlags.saveLeases && v.LeaseID != "" {
				result += fmt.Sprintf("export %s_LEASE_ID=\"%s\"\n", v.Name, v.LeaseID)
			}
		}

		fmt.Print(result)
	},
}
