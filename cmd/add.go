package cmd

import (
	"github.com/Courtcircuits/optique/cli/actions"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new module to the project",
	Long:  `Add a new module to the project`,
	Run: func(cmd *cobra.Command, args []string) {
		actions.AddModule(args[0])
	},
}
