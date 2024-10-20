package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/haokur/dora/cmd"
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var fromFiles []string
var toDir string

var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "选择文件替换对应文件夹下选择要替换的文件",
	Long:  "选择文件替换对应文件夹下选择要替换的文件\n使用：dora replace --to /User/test\n或者：dora replace -f aaa.png -f bbb.png --to /User/test",
	Run: func(cobraCmd *cobra.Command, args []string) {
		workDir := tools.GetWorkDir()
		// 如果用户传入了--target,则使用target，否则列出当前目录下的所有文件供选择
		// 然后列出所有对应得上的目标目录下的文件名，供选择替换
		// 选择后执行替换
		if toDir == "" {
			fmt.Println("请输入要替换的目标目录")
			cobraCmd.Help()
			return
		}
		// 如果from的传入为空，则调用列举当前目录下所有的文件
		if len(fromFiles) == 0 {
			result, err := tools.ReadFilesShallowly(workDir)
			tools.SortFiles(result, "modtime")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			searchChoices := []cmd.CommandItem{}
			for _, fileItem := range result {
				// fileTime := fileItem.LastModified.Format("2006-01-02 15:04:05")
				// fileSize := tools.FormatSize(fileItem.Size)
				searchChoices = append(searchChoices, cmd.CommandItem{
					Value: fileItem.Name,
					Label: "",
					Desc:  "",
				})
			}
			userChoices, err := cmd.Search(searchChoices)
			if err != nil {
				fmt.Println("获取选择要去替换的文件失败", err)
				return
			}
			if len(userChoices) == 0 {
				fmt.Println("未选择要替换的项，自动退出程序")
				os.Exit(1)
			}

			fromFiles = userChoices
		}

		// 筛选出目标目录下文件名匹配的选项
		toDirFiles, err := tools.ReadFilesRecursively(toDir)
		if err != nil {
			fmt.Println("递归读取目标目录文件夹目录失败", err)
			os.Exit(1)
		}
		filterToPaths := []string{}
		for _, item := range toDirFiles {
			for _, choice := range fromFiles {
				if strings.HasSuffix(item, choice) {
					// filterToPaths = append(filterToPaths, fmt.Sprintf("%s/%s", toDir, item))
					filterToPaths = append(filterToPaths, item)
				}
			}
		}
		// 选择要替换的文件
		userSelect2Replace, _, err := cmd.Check("请选择要替换的文件", &filterToPaths, false)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 进行名称匹配的文件替换
		for _, choice := range fromFiles {
			fromFilePath := filepath.Join(workDir, choice)
			for _, fileItem := range userSelect2Replace {
				if strings.HasSuffix(fileItem, choice) {
					toFilePath := filepath.Join(toDir, fileItem)
					// err := os.Rename(fromFilePath, toFilePath)
					err := tools.CopyFile(fromFilePath, toFilePath)
					if err != nil {
						fmt.Println("复制替换文件失败", err)
					} else {
						fmt.Printf("替换 %s 到 %s 成功\n", fromFilePath, toFilePath)
					}
				}
			}
		}
	},
}

func init() {
	replaceCmd.Flags().StringArrayVarP(&fromFiles, "from", "f", []string{}, "输入要替换的文件,可选")
	replaceCmd.Flags().StringVarP(&toDir, "to", "t", "", "输入需要被替换的文件夹，必选")
	rootCmd.AddCommand(replaceCmd)
}
