package cli

import (
	"os"

	prompt "github.com/c-bata/go-prompt"
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

func executor(t string) {
	if t != "" {
		tools.RunCommandWithLog(t)
	}
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "git status", Description: "查看git状态"},
		{Text: "git add .", Description: "添加所有文件改变"},
		{Text: "git commit", Description: "增加git描述"},
		{Text: "git push origin main", Description: "推送提交"},
	}
}

var rootCmd = &cobra.Command{
	Use:   "dora",
	Short: "效率自动化工具箱",
	Long:  `基于Golang+Cobra开发的效率自动化工具箱`,
	Run: func(cmd *cobra.Command, args []string) {
		// 缺省不带参数，则进入dora环境，使用go-prompt进行提示
		p := prompt.New(
			executor,
			completer,
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
