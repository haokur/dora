package cli

import (
	"fmt"
	"strings"

	"github.com/haokur/dora/cmd"
	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var ipv4Flag bool
var ipv6Flag bool
var isCopyFlag bool

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "ip",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ipv4, ipv6 := tools.GetIpAddress()
		ipResult := []string{}
		if ipv4Flag {
			for _, v := range ipv4 {
				if v != "127.0.0.1" {
					ipResult = append(ipResult, v)
				}
			}
		}
		if ipv6Flag {
			ipResult = append(ipResult, ipv6...)
		}

		if isCopyFlag {
			// 复制
			selectIps := ipResult
			// 如果只有一个自动复制
			if len(ipResult) > 1 {
				selectIps, _, _ = cmd.Check("选择要复制的IP地址", &ipResult, false)
			}
			copyStr := strings.Join(selectIps, "\n")
			if copyStr != "" {
				tools.CopyText2ClipBoard(copyStr)
				fmt.Printf("IP地址：%s 已复制到剪切板\n", copyStr)
			}
		} else {
			// 仅打印输出
			for k, v := range ipResult {
				fmt.Println(k+1, v)
			}
		}
	},
}

func init() {
	// dora ip --ipv4=false
	ipCmd.Flags().BoolVarP(&ipv4Flag, "ipv4", "4", true, "是否需要输出ipv4 IP")
	ipCmd.Flags().BoolVarP(&ipv6Flag, "ipv6", "6", false, "是否需要输出ipv6 IP")
	ipCmd.Flags().BoolVarP(&isCopyFlag, "copy", "c", true, "是否需要复制操作")

	rootCmd.AddCommand(ipCmd)
}
