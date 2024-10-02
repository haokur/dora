package cli

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

type noteItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type noteJsonType struct {
	Notes []noteItem `json:"notes"`
}

var noteJsonConfig noteJsonType

func noteExecutor(t string) {
	if t != "exit" {
		tools.CopyText2ClipBoard(t)
		fmt.Println(t, "已复制到剪切板")
	}
	os.Exit(0)
}

func noteCompleter(t prompt.Document) []prompt.Suggest {
	// t.Text中没有空格，则按整条命令来提示
	// t.Text中有空格，则需要将命令按t.Text来匹配再切割，余下的命令字符串
	// 比如无空格，输入gip，能匹配到建议：git push origin main
	// 如果有空格，比如git push，则能匹配到 origin main
	// 如果t.Text为git push origin，则能匹配到main
	noteConfig := noteJsonConfig.Notes
	searchKey := strings.TrimLeft(t.Text, " ")
	suggestions := make([]prompt.Suggest, 0, len(noteConfig))
	matchFieldKey := "Value"
	if tools.ContainsChineseWords(searchKey) {
		matchFieldKey = "Label"
	}
	matches := tools.FindMatches(noteConfig, matchFieldKey, searchKey)

	for _, item := range matches {
		command := item.Value
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

// 备忘笔记本，提供查询列表，可以搜索并复制内容
var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "可搜索复制的备忘命令列表",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tools.ReadDoraJsonConfig(&noteJsonConfig); err != nil {
			fmt.Println("ReadJsonError", err)
			os.Exit(1)
		}

		prefix := "📝notes >>> "

		p := prompt.New(
			noteExecutor,
			noteCompleter,
			prompt.OptionPrefix(prefix),
			prompt.OptionTitle("dora备忘录"),
			prompt.OptionPrefixTextColor(prompt.DarkBlue),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSuggestionBGColor(prompt.DarkGray),
		)
		p.Run()
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
