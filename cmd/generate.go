package cmd

import (
	"github.com/Courtcircuits/optique/cli/actions"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "Bootstrap a new Optique module",
	Long:  `Generates a new Optique module`,
	Run: func(cmd *cobra.Command, args []string) {
		actions.GenerateFromForm(args[0])
	},
}
