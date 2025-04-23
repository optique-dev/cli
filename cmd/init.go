package cmd

import (
	"github.com/Courtcircuits/optique/cli/actions"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "This command will initialize your project",
	Long:  `This command will initialize your project`,
	Run: func(cmd *cobra.Command, args []string) {
		init := actions.NewInitialization(args[0])
		actions.Initialize(init)
	},
}
