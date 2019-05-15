package command

import (
	"fmt"
	"os"
)

// Run is used to start command-line interface
func Run() {
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(renewCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
