package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/haokur/dora/tools"
	"github.com/spf13/cobra"
)

var configPath string

var updateFlag bool
var infoFlag bool

// 从远程下载
var downloadKey string

// 更新到远程
var publishFlag bool

type configRespData struct {
	Api_key string `json:"api_key"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
type configResp struct {
	Code int            `json:"code"`
	Data configRespData `json:"data"`
	Msg  string         `json:"msg"`
}

type doraConfigType struct {
	Api_key string
	Name    string
}

// 下载配置
func downloadConfig() error {
	remoteUrl := fmt.Sprintf("http://106.53.114.178:8008/dora_config/item?api_key=%s", downloadKey)
	resp, err := http.Get(remoteUrl)
	if err != nil {
		fmt.Println("无法读取远程配置，使用空配置初始化")
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应的 JSON 数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容时出错: %v", err)
	}

	configData := configResp{}
	err = json.Unmarshal(body, &configData)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	content := configData.Data.Content
	if content != "" {
		err = os.WriteFile(tools.GetDoraConfigPath(), []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("写入本地配置文件时出错: %v", err)
		}
	}

	fmt.Println("下载配置文件成功", content)

	return nil
}

// 上传配置
type Payload struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func publishConfig() error {
	// 读取本地配置
	localConfig := doraConfigType{}
	jsonFilePath := tools.GetDoraConfigPath()

	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &localConfig); err != nil {
		return err
	}

	api_key := localConfig.Api_key
	// remoteUrl := fmt.Sprintf("http://localhost:8008/dora_config/update?api_key=%s", api_key)
	remoteUrl := fmt.Sprintf("http://106.53.114.178:8008/dora_config/update?api_key=%s", api_key)

	// 创建负载
	payload := Payload{
		Name:    localConfig.Name,
		Content: string(data),
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(remoteUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码：%d", resp.StatusCode)
	}

	return nil
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理配置文件，位于用户目录/dora/.config.json",
	Run: func(cmd *cobra.Command, args []string) {
		// 如果是查看配置文件
		if infoFlag {
			tools.PreviewFileWithSystemEditor(configPath)
			return
		}

		// 更新配置
		if updateFlag {
			tools.EditFileWithSystemEditor(configPath)
			return
		}

		// 从远程拉取配置
		// dora config -d [api_key]
		if downloadKey != "" {
			downloadConfig()
			return
		}

		if publishFlag {
			publishConfig()
			return
		}

		cmd.Help()
	},
}

func init() {
	userHomeDir, _ := os.UserHomeDir()
	configPath = filepath.Join(userHomeDir, "dora/.config.json")

	configCmd.Flags().BoolVarP(&infoFlag, "info", "i", false, "查看配置信息")
	configCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "更新配置文件")
	configCmd.Flags().StringVarP(&downloadKey, "download", "d", "", "从远程拉取配置，dora config -d [api_key]")
	configCmd.Flags().BoolVarP(&publishFlag, "publish", "p", false, "将配置推送到远程,使用配置中的api_key，dora config -d")
	rootCmd.AddCommand(configCmd)
}
