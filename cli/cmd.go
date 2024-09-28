package cli

import (
	"fmt"
	"os"

	"github.com/haokur/dora/cmd"
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

type cmdJsonItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type cmdJsonType struct {
	Commands []cmdJsonItem `json:"commands"`
}

var cmdTip = &cobra.Command{
	Use:   "cmd",
	Short: "命令行提示,自动输入",
	Run: func(cobraCmd *cobra.Command, args []string) {
		jsonFilePath := "./configs/cmd.json"

		var jsonData cmdJsonType
		if err := tools.ReadJsonFile(jsonFilePath, &jsonData); err != nil {
			fmt.Println("ReadJsonError", err)
			os.Exit(1)
		}

		// 类型转化
		searchParams := tools.Convert(jsonData.Commands, func(cmdItem cmdJsonItem) cmd.CommandItem {
			return cmd.CommandItem{
				Label: cmdItem.Label,
				Value: cmdItem.Value,
			}
		})

		result, err := cmd.Search(searchParams)
		if err != nil {
			fmt.Println("cmd Search error", err)
		}
		for _, v := range result {
			err := tools.RunCommandWithLog(v)
			if err != nil {
				fmt.Println("执行失败", v, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cmdTip)
}
