package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/stackdumper/vault-to-env/client"
)

func getByteIndex(s string, b byte) int {
	index := strings.IndexByte(s, b)

	if index == -1 {
		return len(s)
	}

	return index
}

type Var struct {
	Var, Path, Key, Value string
}

// Run is used to start command-line interface
func Run() {
	// initialize client
	client, err := client.NewClient(&client.Config{
		VaultAddress: "http://localhost:8200",
		VaultToken:   "s.QqzFRLhsMPcNxIBu6ZE34pMV",
	})
	if err != nil {
		log.Fatal(err)
	}

	var options struct {
		output       string
		envs         FlagList
		vaultAddress string
		vaultToken   string
	}

	flag.Var(&options.envs, "e", "environment variables to fetch")
	flag.StringVar(&options.output, "o", "", "output file path")
	flag.StringVar(&options.vaultAddress, "a", "http://localhost:8200", "vault address")
	flag.StringVar(&options.vaultToken, "t", "", "vault token")
	flag.Parse()

	// validate flags
	if len(options.envs) == 0 {
		flag.Usage()
		log.Fatal("-e flag is required")
	}

	// prepare args
	var vars = make([]Var, len(options.envs))
	for index, env := range options.envs {
		iv, ip, ik := getByteIndex(env, '='), getByteIndex(env, '#'), len(env)

		vars[index].Var = env[:iv]
		vars[index].Path = env[iv+1 : ip]
		vars[index].Key = env[ip+1 : ik]
	}

	// fetch values
	for i, v := range vars {
		result := client.Read(v.Path, v.Key)

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

	if options.output == "" {
		// write to stdout if output file path is not provided

		fmt.Print(result)
	} else {
		// write to file otherwise

		err = ioutil.WriteFile(options.output, []byte(result), 0740)
		if err != nil {
			log.Fatal(err)
		}
	}
}
