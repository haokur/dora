package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// 读取 JSON 文件并将数据解码为指定的类型
func ReadJsonFile[T any](filePath string, out *T) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, out); err != nil {
		return err
	}
	return nil
}

// 读取dora的json配置
func ReadDoraJsonConfig[T any](out *T) error {
	userHomeDir, _ := os.UserHomeDir()
	jsonFilePath := filepath.Join(userHomeDir, "dora/.config.json")

	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, out); err != nil {
		return err
	}
	return nil
}
