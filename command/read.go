package command

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var saveFlags struct {
	vars       []string
	saveLeases bool
}

func init() {
	saveCmd.Flags().StringSliceVar(&saveFlags.vars, "vars", []string{}, "list of vars to read")
	saveCmd.Flags().BoolVar(&saveFlags.saveLeases, "save-leases", false, "save secret leases")
}

type ResultVar struct {
	// Input: "Name=Path#Key"
	// Output: "export Name=Value"
	Name, Path, Key, Value, Lease string
}

// Read secrets and output them as env variables
var saveCmd = &cobra.Command{
	Use:   "read",
	Short: "Read secrets and output them as env variables",
	Run: func(cmd *cobra.Command, args []string) {
		var resultVars = make([]ResultVar, len(saveFlags.vars))

		for index, rawVar := range saveFlags.vars {
			// Var=Path#Key
			iv, ip, ik := getByteIndex(rawVar, '='), getByteIndex(rawVar, '#'), len(rawVar)

			// exit if provided var is incompliete
			if len(rawVar) <= ip+1 {
				log.Fatal("unrecognized variable format, required var=path#key")
			}

			resultVars[index] = ResultVar{
				Name: rawVar[:iv],
				Path: rawVar[iv+1 : ip],
				Key:  rawVar[ip+1 : ik],
			}
		}

		// get Values
		for i, v := range resultVars {
			result := PersistentState.client.Read(v.Path, v.Key)

			// print warnings if there are some
			if len(result.Warnings) != 0 {
				for _, warning := range result.Warnings {
					log.Println(warning)
				}
			}

			// exit with error if there's one
			if result.Error != nil {
				log.Fatal(result.Error)
			}

			// assign value
			resultVars[i].Value = result.Value
			resultVars[i].Lease = result.LeaseID
		}

		// extract result string
		// export Var="Value"
		var result string
		for _, v := range resultVars {
			result += fmt.Sprintf("export %s=\"%s\"\n", v.Name, v.Value)

			if saveFlags.saveLeases && v.Lease != "" {
				result += fmt.Sprintf("export %s_LEASE=\"%s\"\n", v.Name, v.Lease)
			}
		}

		fmt.Print(result)
	},
}
