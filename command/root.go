package command

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
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
		if (PersistentFlags.AuthPath == "" || len(PersistentFlags.AuthData) == 0) && PersistentFlags.AuthToken == "" {
			log.Fatal("either auth-path and auth-data or auth-token must be provided")
		}

		// parse auth data into map[string]interface{}
		var authData = make(map[string]interface{})
		for _, v := range PersistentFlags.AuthData {
			splitted := strings.Split(v, "=")

			if len(splitted) != 2 {
				log.Fatal("malformed auth data format (required key=value)")
			}

			authData[splitted[0]] = splitted[1]
		}

		// initialize the client
		client, err := client.NewClient(&client.Config{
			VaultAddress: PersistentFlags.Address,
		})
		if err != nil {
			log.Fatal(err)
		}

		// authenticate
		token, err := client.Auth(PersistentFlags.AuthPath, authData)
		if err != nil {
			log.Fatal(err)
		}

		// set token
		client.SetToken(token)

		// save client
		PersistentState.client = client
	},
}
