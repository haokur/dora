package tools

import (
	"fmt"
	"strings"
)

// isSubsequence 检查短字符串是否为长字符串的子序列
func IsSubsequence(short, long string) bool {
	shortLen, longLen := len(short), len(long)
	if shortLen == 0 {
		return true // 空字符串是任何字符串的子序列
	}

	j := 0 // 指针指向短字符串
	for i := 0; i < longLen; i++ {
		if long[i] == short[j] {
			j++
			if j == shortLen {
				return true // 完全匹配
			}
		}
	}
	return false // 没有完全匹配
}

// 获取对像数组里按某个key，以IsSubsequence来查找匹配的值
func FindMatches(arr []map[string]interface{}, fieldKey string, searchKey string) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, item := range arr {
		if value, ok := item[fieldKey].(string); ok {
			if IsSubsequence(searchKey, value) {
				result = append(result, item)
			}
		}
	}
	return result
}

// 获取高亮匹配的字符串
func GetHighlightString(command, input string) string {
	if input == "" {
		return command // 如果没有输入，直接返回命令
	}

	var result strings.Builder
	inputLen := len(input)
	commandLen := len(command)

	// 记录输入在命令中的位置
	inputIndex := 0
	for i := 0; i < commandLen; i++ {
		if inputIndex < inputLen && command[i] == input[inputIndex] {
			result.WriteString(fmt.Sprintf("\033[1;31;4m%s\033[0m", string(command[i]))) // 高亮并下划线
			inputIndex++
		} else {
			result.WriteString(string(command[i])) // 正常显示
		}
	}

	return result.String()
}
