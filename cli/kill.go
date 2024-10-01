package cli

import (
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill",
	Run: func(cmd *cobra.Command, args []string) {
		tools.KillProcess(&args)
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
