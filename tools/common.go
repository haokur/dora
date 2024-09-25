package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 获取用户的路径
func GetUserHomePath() string {
	dirPath, err := os.UserHomeDir()
	if err != nil {
		return "~"
	}
	return dirPath
}

// 获取git根目录
func GetGitRootDir() (string, error) {
	// 使用 'git rev-parse --show-toplevel' 获取Git根目录
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git root directory: %w", err)
	}

	// 去除输出中的换行符和空白
	gitRootDir := strings.TrimSpace(string(out))
	return gitRootDir, nil
}

func readDirRecursively(dirPath string, filePaths *[]string, rootPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// 遍历当前目录中的所有条目
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		// 计算相对于 rootPath 的路径
		relativePath, _ := filepath.Rel(rootPath, fullPath)

		// 如果是目录，递归处理
		if entry.IsDir() {
			err := readDirRecursively(fullPath, filePaths, rootPath)
			if err != nil {
				return err
			}
		} else {
			*filePaths = append(*filePaths, relativePath)
		}
	}
	return nil
}

// 递归读取文件函数
func ReadFilesRecursively(dirPath string) ([]string, error) {
	result := []string{}
	err := readDirRecursively(dirPath, &result, dirPath)
	return result, err
}
