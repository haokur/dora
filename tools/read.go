package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

// 读取json文件
func ReadJsonFile(filePath string) (map[string]interface{}, error) {
	var config map[string]interface{}
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(bytes, &config)
	return config, err
}

// 读取json的字段值，且其值为字符串数组
func ReadJsonFieldValueAsSlice(filePath string, field string) ([]string, error) {
	result, err := ReadJsonFile(filePath)
	if err != nil {
		fmt.Println("ReadJsonError", err)
		return nil, err
	}

	strings := []string{}
	if stringInterfaces, ok := result[field].([]interface{}); ok {
		for _, cmd := range stringInterfaces {
			if strCmd, ok := cmd.(string); ok {
				strings = append(strings, strCmd)
			}
		}
	}
	return strings, nil
}
