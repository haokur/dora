package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "监听文件变化自动执行命令",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("watch is running")
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
