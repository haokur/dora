package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

// 返回拼接好的目录树字符串
func printTree(dir string, indent string, ignoredDirs []string) (string, error) {
	var result string
	// 打开目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// 先处理目录
	for _, entry := range entries {
		if entry.IsDir() {
			// 检查当前条目是否在忽略列表中
			if tools.SliceContains(ignoredDirs, entry.Name()) {
				continue
			}
			// 拼接目录条目
			result += indent + entry.Name() + "/\n"
			nextDir := filepath.Join(dir, entry.Name())
			subtree, err := printTree(nextDir, "   "+indent, ignoredDirs)
			if err != nil {
				return "", err
			}
			result += subtree
		}
	}

	// 再处理文件
	for _, entry := range entries {
		if !entry.IsDir() {
			// 拼接文件条目
			result += indent + entry.Name() + "\n"
		}
	}

	return result, nil
}

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "递归打印文件夹文件树状结构",
	Run: func(cmd *cobra.Command, args []string) {
		// 读取当前目录
		workDir := tools.GetWorkDir()
		// 要忽略的目录
		ignoredDirs := []string{
			".git",
			"node_modules",
			"dist",
		}
		str, err := printTree(workDir, "- ", ignoredDirs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(str)
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
}
