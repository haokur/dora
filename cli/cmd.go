package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdTip = &cobra.Command{
	Use:   "cmd",
	Short: "命令行提示,自动输入",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func init() {
	rootCmd.AddCommand(cmdTip)
}
