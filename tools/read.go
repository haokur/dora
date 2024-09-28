package tools

import (
	"encoding/json"
	"os"
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
