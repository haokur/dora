package cli

import (
	"fmt"

	"github.com/haokur/dora/cmd"
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var cmdTip = &cobra.Command{
	Use:   "cmd",
	Short: "命令行提示,自动输入",
	Run: func(cobraCmd *cobra.Command, args []string) {
		jsonFilePath := "./configs/cmd.json"
		commands, err := tools.ReadJsonFieldValueAsSlice(jsonFilePath, "commands")
		if err != nil {
			fmt.Println("ReadJsonFieldValueAsSlice error:", err)
			return
		}
		result, err := cmd.Search(commands)
		if err != nil {
			fmt.Println("cmd Search error", err)
		}
		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(cmdTip)
}
