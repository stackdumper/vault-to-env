package command

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackdumper/vault-to-env/client"
)

var PersistentState struct {
	client *client.Client
}

var PersistentFlags struct {
	Address  string
	AuthPath string
	AuthData []string
}

func init() {
	rootCmd.PersistentFlags().StringVar(&PersistentFlags.Address, "address", "http://localhost:8200", "Vault address")
	rootCmd.PersistentFlags().StringVar(&PersistentFlags.AuthPath, "auth-path", "", "Vault auth path")
	rootCmd.PersistentFlags().StringSliceVar(&PersistentFlags.AuthData, "auth-data", []string{}, "Vault auth data")

	rootCmd.MarkFlagRequired("auth-path")
	rootCmd.MarkFlagRequired("auth-data")
}

// root command
var rootCmd = &cobra.Command{
	Use:   "vte",
	Short: "A tool to manipulate for saving and manipulating secrets from Hashicorp Vault",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
