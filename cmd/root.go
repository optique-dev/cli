package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use: "optique",
	Short: "Focus only on your business logic",
	Long: `Focus only on your business logic`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
