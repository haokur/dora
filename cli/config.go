package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var configPath string

var updateFlag bool
var infoFlag bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理配置文件",
	Run: func(cmd *cobra.Command, args []string) {
		// 如果是查看配置文件
		if infoFlag {
			tools.RunCommandWithLog(fmt.Sprintf("cat %s", configPath))
			return
		}

		// 更新配置
		if updateFlag {
			tools.RunCommandWithLog(fmt.Sprintf("code %s", configPath))
			return
		}
	},
}

func init() {
	userHomeDir, _ := os.UserHomeDir()
	configPath = filepath.Join(userHomeDir, "dora/.config.json")

	configCmd.Flags().BoolVarP(&infoFlag, "info", "i", false, "查看配置信息")
	configCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "更新配置文件")
	rootCmd.AddCommand(configCmd)
}
