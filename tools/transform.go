package tools

// 通用转换函数
func Convert[T any, R any](source []T, convertFunc func(T) R) []R {
	var result []R
	for _, item := range source {
		result = append(result, convertFunc(item))
	}
	return result
}
