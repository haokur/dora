package tools

import (
	"fmt"
	"reflect"
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

func FindMatches[T any](arr []T, fieldKey string, searchKey string) []T {
	result := []T{}
	for _, item := range arr {
		// 使用反射来获取字段值
		val := reflect.ValueOf(item).FieldByName(fieldKey)
		if val.IsValid() && val.Kind() == reflect.String {
			if IsSubsequence(searchKey, val.String()) {
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
