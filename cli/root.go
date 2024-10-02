package cli

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

func getPrefix() string {
	workDir := tools.GetWorkDir()
	homeDir := tools.GetUserHomePath()
	shortWorkDir := strings.ReplaceAll(workDir, homeDir, "~")
	prefix := fmt.Sprintf("⚡%s >>> ", shortWorkDir)
	return prefix
}

func createPrompt() *prompt.Prompt {
	prefix := getPrefix()
	return prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(prefix),
		prompt.OptionTitle("dora命令行工具"),
		prompt.OptionPrefixTextColor(prompt.DarkBlue),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	)
}

func executor(t string) {
	if t == "dora" {
		return
	}
	if t == "exit" {
		fmt.Println("再见！")
		os.Exit(0)
	}
	if t != "" {
		err := tools.RunCommandWithLog(t)
		if err == nil && strings.HasPrefix(t, "cd") {
			p := createPrompt()
			p.Run()
		}
	}
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
	searchKey := strings.TrimLeft(t.Text, " ")
	suggestions := make([]prompt.Suggest, 0, len(promptConfig))
	matchFieldKey := "Cmd"
	if tools.ContainsChineseWords(searchKey) {
		matchFieldKey = "Label"
	}
	matches := tools.FindMatches(promptConfig, matchFieldKey, searchKey)

	for _, item := range matches {
		command := item.Cmd
		if strings.Contains(searchKey, " ") {
			// 替换最后一个空格前面所有内容
			beforeCmd := tools.GetBeforeLastSpace(searchKey) + " "
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
	Long:  "基于Golang+Cobra开发的效率自动化工具箱\n不带参数进入带提示的交互页面",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tools.ReadDoraJsonConfig(&jsonConfig); err != nil {
			fmt.Println("ReadJsonError", err)
			os.Exit(1)
		}

		// 缺省不带参数，则进入dora环境，使用go-prompt进行提示
		p := createPrompt()
		p.Run()

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
