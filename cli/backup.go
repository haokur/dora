package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var isRecover bool
var isBackup bool
var isWithOpen bool
var backupFileName string

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "备份git未提交的代码，备份目录/Users/dora/backup/项目名_备份日期",
	Run: func(cmd *cobra.Command, args []string) {
		userHomeDir, _ := os.UserHomeDir()
		gitBackupBaseDir := fmt.Sprintf("%s/%s", userHomeDir, "dora/backup")
		currentWorkGitDir, err := tools.GetGitRootDir()
		if err != nil {
			fmt.Println("获取git根目录失败", err)
			return
		}

		fileName := backupFileName
		if fileName == "" {
			fileName = filepath.Base(currentWorkGitDir)
		}

		gitBackupDir := fmt.Sprintf("%s/%s", gitBackupBaseDir, fileName)
		if isBackup {
			if err != nil {
				return
			}
			backupPath, err := tools.BackupUnCommitFiles(currentWorkGitDir, gitBackupDir)
			if err != nil {
				fmt.Println("备份失败:", err)
				return
			}
			if isWithOpen {
				tools.OpenFolderAndSelectFile(backupPath)
			}
		} else if isRecover {
			// 1.找到匹配的备份目录
			// 2.以时间戳按时间倒序，最近的备份显示在最前面，单选
			// 3.用户选择一个备份目录，点击确认
			// 4.展示选择备份目录下所有文件，且显示更改时间，文件大小，用户选择要还原的文件
			// 5.将用户选择的文件，还原到git项目目录
			tools.RecoverBackupFiles(gitBackupBaseDir, currentWorkGitDir)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	backupCmd.Flags().BoolVarP(&isBackup, "backup", "b", false, "备份文件")
	backupCmd.Flags().BoolVarP(&isRecover, "cover", "c", false, "恢复文件")
	backupCmd.Flags().BoolVarP(&isWithOpen, "open", "o", false, "完成后打开")
	backupCmd.Flags().StringVarP(&backupFileName, "name", "n", "", "备份文件名")
	rootCmd.AddCommand(backupCmd)
}
