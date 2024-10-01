package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "备份git未提交的代码",
	Run: func(cmd *cobra.Command, args []string) {
		userHomeDir, _ := os.UserHomeDir()
		currentWorkGitDir, err := tools.GetGitRootDir()
		if err != nil {
			fmt.Println("获取git根目录失败", err)
			return
		}
		action := args[0]
		switch action {
		// 恢复
		case "recover":
			// 1.找到匹配的备份目录
			// 2.以时间戳按时间倒序，最近的备份显示在最前面，单选
			// 3.用户选择一个备份目录，点击确认
			// 4.展示选择备份目录下所有文件，且显示更改时间，文件大小，用户选择要还原的文件
			// 5.将用户选择的文件，还原到git项目目录
			fileName := filepath.Base(currentWorkGitDir)
			gitBackupDir := fmt.Sprintf("%s/%s", userHomeDir, fileName)
			tools.RecoverBackupFiles(gitBackupDir, currentWorkGitDir)
		case "backup":
			fileName := filepath.Base(currentWorkGitDir)
			if err != nil {
				return
			}
			backupTargetDir := fmt.Sprintf("%s/%s", userHomeDir, fileName)

			tools.BackupUnCommitFiles(currentWorkGitDir, backupTargetDir)
		default:
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
