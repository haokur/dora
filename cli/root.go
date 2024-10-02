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
	Cmd      string       `json:"cmd"`
	Label    string       `json:"label"`
	Children []promptItem `json:"children"`
}

type promptJsonType struct {
	Prompts []promptItem `json:"prompts"`
}

var jsonConfig promptJsonType

func getSuggestions(input string, commands []promptItem) []prompt.Suggest {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	currentCommands := commands
	for _, part := range parts {
		found := false
		for _, cmd := range currentCommands {
			if cmd.Cmd == part {
				currentCommands = cmd.Children
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}

	var suggestions []prompt.Suggest
	for _, cmd := range currentCommands {
		suggestions = append(suggestions, prompt.Suggest{Text: cmd.Cmd, Description: cmd.Label})
	}
	return suggestions
}

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

	// beforeCmd为以空格切割的命令
	beforeCmd := tools.GetBeforeLastSpace(searchKey)
	searchKeyHasSpace := strings.Contains(searchKey, " ")

	// 使用beforeCmd，先全等比较，找到对应下的children
	// 再使用空格切开，逐级拼接查找，找到对应下的children
	if searchKeyHasSpace {
		// 空格前的筛选接下来的列表，空格后的在筛下来的列表里继续筛
		subSuggestList := getSuggestions(beforeCmd, promptConfig)
		afterLastSpaceCmd := strings.ReplaceAll(searchKey, beforeCmd+" ", "")
		afterLastSpaceCmd = strings.Trim(afterLastSpaceCmd, "")
		filterChildrenCmds := subSuggestList
		if afterLastSpaceCmd != "" {
			filterChildrenCmds = tools.FindMatches(subSuggestList, "Text", afterLastSpaceCmd)
		}
		suggestions = append(suggestions, filterChildrenCmds...)
	}

	matches := tools.FindMatches(promptConfig, matchFieldKey, searchKey)
	for _, item := range matches {
		command := item.Cmd
		if searchKeyHasSpace {
			// 替换最后一个空格前面所有内容
			command = strings.ReplaceAll(command, beforeCmd+" ", "")
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
