package command

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	cobra "github.com/spf13/cobra"

	"github.com/stackdumper/vault-to-env/client"
)

var PersistentState struct {
	client *client.Client
}

var PersistentFlags struct {
	Address   string
	AuthPath  string
	AuthData  []string
	AuthToken string
}

func init() {
	rootCmd.PersistentFlags().StringVar(&PersistentFlags.Address, "address", os.Getenv("VAULT_ADDR"), "Vault address")
	rootCmd.PersistentFlags().StringVar(&PersistentFlags.AuthPath, "auth-path", "", "Vault auth path")
	rootCmd.PersistentFlags().StringSliceVar(&PersistentFlags.AuthData, "auth-data", []string{}, "Vault auth data")
	rootCmd.PersistentFlags().StringVar(&PersistentFlags.AuthToken, "auth-token", os.Getenv("VAULT_TOKEN"), "Vault auth token")
}

// root command
var rootCmd = &cobra.Command{
	Use:   "vte",
	Short: "A tool to manipulate for saving and manipulating secrets from Hashicorp Vault",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// setup logger
		if os.Getenv("DEBUG") == "1" {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}

		if (PersistentFlags.AuthPath == "" || len(PersistentFlags.AuthData) == 0) && PersistentFlags.AuthToken == "" {
			log.WithFields(log.Fields{
				"auth-path":  PersistentFlags.AuthPath,
				"auth-data":  PersistentFlags.AuthData,
				"auth-token": PersistentFlags.AuthToken,
			}).Fatal("either auth-path and auth-data or auth-token must be provided")
		}

		// parse auth data into map[string]interface{}
		log.WithField("data", PersistentFlags.AuthData).Debug("parsing auth data")
		var authData = make(map[string]interface{})
		for _, v := range PersistentFlags.AuthData {
			splitted := strings.Split(v, "=")

			if len(splitted) != 2 {
				log.WithFields(log.Fields{
					"data":     v,
					"splitted": splitted,
				}).Fatal("unrecognized auth data format (required key=value)")
			}

			authData[splitted[0]] = splitted[1]
		}

		// initialize the client
		log.WithField("address", PersistentFlags.Address).Debug("initializing vault client")
		client, err := client.NewClient(&client.Config{
			VaultAddress: PersistentFlags.Address,
		})
		if err != nil {
			log.WithError(err).Fatal("failed to initialize vault client")
		}

		// auth with token
		// if there's no token, get token via auth call
		var token = PersistentFlags.AuthToken
		if token == "" {
			log.WithFields(log.Fields{
				"path": PersistentFlags.AuthPath,
				"data": authData,
			}).Debug("authenticating")
			token, err = client.Auth(PersistentFlags.AuthPath, authData)
			if err != nil {
				log.WithError(err).Fatal("failed to authenticate vault client")
			}
		}

		// save tokem
		log.WithField("token", token).Debug("saving token")
		client.SetToken(token)

		// save client
		PersistentState.client = client
	},
}
