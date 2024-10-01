package cli

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

func executor(t string) {
	if t == "dora" {
		return
	}
	if t == "exit" {
		fmt.Println("再见！")
		os.Exit(0)
	}
	if t != "" {
		tools.RunCommandWithLog(t)
	}
}

func getBeforeLastSpace(input string) string {
	// 查找最后一个空格的位置
	lastSpaceIndex := strings.LastIndex(input, " ")
	if lastSpaceIndex == -1 {
		// 如果没有空格，返回原字符串
		return input
	}
	// 返回空格前的部分
	return input[:lastSpaceIndex]
}

type promptItem struct {
	Cmd   string `json:"cmd"`
	Label string `json:"label"`
}

type promptJsonType struct {
	Prompts []promptItem `json:"prompts"`
}

var jsonConfig promptJsonType

func completer(t prompt.Document) []prompt.Suggest {
	// t.Text中没有空格，则按整条命令来提示
	// t.Text中有空格，则需要将命令按t.Text来匹配再切割，余下的命令字符串
	// 比如无空格，输入gip，能匹配到建议：git push origin main
	// 如果有空格，比如git push，则能匹配到 origin main
	// 如果t.Text为git push origin，则能匹配到main
	promptConfig := jsonConfig.Prompts
	searchKey := t.Text
	suggestions := make([]prompt.Suggest, 0, len(promptConfig))
	matches := tools.FindMatches(promptConfig, "Cmd", searchKey)

	for _, item := range matches {
		command := item.Cmd
		if strings.Contains(searchKey, " ") {
			// 替换最后一个空格前面所有内容
			beforeCmd := getBeforeLastSpace(searchKey) + " "
			command = strings.ReplaceAll(command, beforeCmd, "")
		}
		suggestions = append(suggestions, prompt.Suggest{
			Text:        command,
			Description: item.Label,
		})
	}

	return suggestions
}

var rootCmd = &cobra.Command{
	Use:   "dora",
	Short: "效率自动化工具箱",
	Long:  `基于Golang+Cobra开发的效率自动化工具箱`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tools.ReadDoraJsonConfig(&jsonConfig); err != nil {
			fmt.Println("ReadJsonError", err)
			os.Exit(1)
		}

		// 缺省不带参数，则进入dora环境，使用go-prompt进行提示
		p := prompt.New(
			executor,
			completer,
			prompt.OptionPrefix("⚡dora >>> "),
			prompt.OptionTitle("dora命令行工具"),
			// prompt.OptionHistory([]string{"SELECT * FROM users;"}), // 设置初始化历史记录可上下翻动
			prompt.OptionPrefixTextColor(prompt.DarkBlue),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSuggestionBGColor(prompt.DarkGray),
		)
		p.Run()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
