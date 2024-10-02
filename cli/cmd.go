package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/haokur/dora/cmd"
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

type cmdJsonItem struct {
	Value    string        `json:"value"`
	Label    string        `json:"label"`
	Children []cmdJsonItem `json:"children"`
}

type cmdJsonType struct {
	Commands []cmdJsonItem `json:"commands"`
}

var cmdTip = &cobra.Command{
	Use:   "cmd",
	Short: "列举dora配置文件中的所有命令，可筛选多选命令依次执行",
	Run: func(cobraCmd *cobra.Command, args []string) {
		// jsonFilePath := "./configs/cmd.json"
		// userHomeDir, _ := os.UserHomeDir()
		// jsonFilePath := filepath.Join(userHomeDir, "dora/.config.json")
		var jsonData cmdJsonType
		if err := tools.ReadDoraJsonConfig(&jsonData); err != nil {
			fmt.Println("ReadJsonError", err)
			os.Exit(1)
		}

		// 类型转化
		searchParams := tools.Convert(jsonData.Commands, func(cmdItem cmdJsonItem) cmd.CommandItem {
			childCmds := []string{}
			if len(cmdItem.Children) > 0 {
				for _, item := range cmdItem.Children {
					childCmds = append(childCmds, item.Value)
				}
			}
			return cmd.CommandItem{
				Label: cmdItem.Label,
				Value: cmdItem.Value,
				Desc:  strings.Join(childCmds, "，"),
			}
		})

		result, err := cmd.Search(searchParams)
		if err != nil {
			fmt.Println("cmd Search error", err)
		}
		for _, v := range result {
			waitRunCmds := []string{}
			// 找到匹配的命令，假如有children属性，执行children里面的内容
			for _, item := range jsonData.Commands {
				if item.Value == v {
					if len(item.Children) > 0 {
						for _, child := range item.Children {
							waitRunCmds = append(waitRunCmds, child.Value)
						}
					} else {
						waitRunCmds = append(waitRunCmds, v)
					}
				}
			}
			for _, cmdItem := range waitRunCmds {
				err := tools.RunCommandWithLog(cmdItem)
				if err != nil {
					fmt.Println("执行失败", v, err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cmdTip)
}
