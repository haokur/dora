package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/haokur/dora/cli"
	"github.com/haokur/dora/tools"
)

// configExists 检查配置文件是否存在
func configExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// 从远程 URL 下载 JSON 文件并生成本地配置文件
func downloadConfig(url string, path string) error {
	// 发送 HTTP GET 请求获取远程 JSON 文件
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("无法读取远程配置，使用空配置初始化")
		defaultConfig := []byte("default_config_key: default_value\n")
		// err := os.WriteFile(path, defaultConfig, 0644)
		tools.SafeWriteFile(path, defaultConfig)
		return err
		// return fmt.Errorf("无法获取远程配置文件: %v", err)
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

	// 将响应的 JSON 数据写入本地文件
	err = os.WriteFile(path, body, 0644)
	if err != nil {
		return fmt.Errorf("写入本地配置文件时出错: %v", err)
	}

	return nil
}

func initConfigAuto() {
	userHomeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(userHomeDir, "dora/.config.json")
	configRemoteUrl := "https://gitee.com/haokur/public-configs/releases/download/dora1.0/dora.config.json"
	if !configExists(configPath) {
		err := downloadConfig(configRemoteUrl, configPath)
		if err != nil {
			fmt.Println("生成配置文件时出错:", err)
		} else {
			fmt.Println("配置文件生成成功:", configPath)
		}
	}
}

func main() {
	initConfigAuto()
	cli.Execute()
}
