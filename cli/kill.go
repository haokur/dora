package cli

import (
	"fmt"

	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var silenceFlag bool

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "清理端口或进程，可同时多个，kill 5173 node nginx",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("请输入要清理的端口或程序名（可多个）,如dora kill 5173 或dora kill 5173 node")
			return
		}
		tools.KillProcess(&args, silenceFlag)
	},
}

func init() {
	killCmd.Flags().BoolVarP(&silenceFlag, "silence", "s", false, "静默清理（无选择步骤）")
	rootCmd.AddCommand(killCmd)
}
