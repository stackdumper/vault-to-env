package command

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

var saveFlags struct {
	vars []string
	out  string
}

func init() {
	saveCmd.Flags().StringSliceVar(&saveFlags.vars, "vars", []string{}, "")
	saveCmd.Flags().StringVar(&saveFlags.out, "out", "", "")
}

type Var struct {
	Var, Path, Key, Value string
}

// Read secrets and output them as env variables
var saveCmd = &cobra.Command{
	Use:   "read",
	Short: "Read secrets and output them as env variables",
	Run: func(cmd *cobra.Command, args []string) {
		// prepare args
		var vars = make([]Var, len(saveFlags.vars))
		for index, env := range saveFlags.vars {
			iv, ip, ik := getByteIndex(env, '='), getByteIndex(env, '#'), len(env)

			vars[index].Var = env[:iv]
			vars[index].Path = env[iv+1 : ip]
			vars[index].Key = env[ip+1 : ik]
		}

		// fetch values
		for i, v := range vars {
			result := PersistentState.client.Read(v.Path, v.Key)

			if len(result.Warnings) != 0 {
				for _, warning := range result.Warnings {
					log.Println(warning)
				}
			}

			if result.Error != nil {
				log.Fatal(result.Error)
			}

			vars[i].Value = result.Value
		}

		// prepare result string
		var result string
		for _, v := range vars {
			result += fmt.Sprintf("export %s=\"%s\"\n", v.Var, v.Value)
		}

		if saveFlags.out == "" {
			// write to stdout if output file path is not provided

			fmt.Print(result)
		} else {
			// write to file otherwise

			err := ioutil.WriteFile(saveFlags.out, []byte(result), 0740)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
