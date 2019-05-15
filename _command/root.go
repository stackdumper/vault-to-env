package command

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "vte",
	Run: func(cmd *cobra.Command, args []string) {},
}
