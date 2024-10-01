package cli

import (
	"fmt"

	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "杀进程",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("请输入要清理的端口或程序名（可多个）,如dora kill 5173 或dora kill 5173 node")
			return
		}
		tools.KillProcess(&args)
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
}
