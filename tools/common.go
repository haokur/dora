package tools

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/atotto/clipboard"
	terminal "golang.org/x/term"
)

// 获取用户的路径
func GetUserHomePath() string {
	dirPath, err := os.UserHomeDir()
	if err != nil {
		return "~"
	}
	return dirPath
}

// 获取工作目录
func GetWorkDir() string {
	workingDir, _ := os.Getwd()
	return workingDir
}

// 安全创建文件，避免文件夹不存在的情况
func SafeWriteFile(filePath string, content []byte) {
	dirPath := path.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	os.WriteFile(filePath, content, 0644)
}

// 是否是要调用终端的vim
func isCallTerminalVim(command string) bool {
	parts := strings.Fields(command)
	isGitCommit := false
	isVim := false
	if len(parts) > 1 {
		isGitCommit = parts[0] == "git" && parts[1] == "commit"
	} else if len(parts) > 0 {
		isVim = strings.HasPrefix(command, "vi")
	}
	return isGitCommit || isVim
}

// 调用系统的vim
func callTerminalVim(command string) {
	// 获取当前终端
	fd := int(os.Stdin.Fd())

	// 设置终端为原始模式
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	// 调用 Vim
	parts := strings.Fields(command) // 使用 Fields 分割以处理空格
	process := parts[0]
	args := parts[1:]
	cmd := exec.Command(process, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		panic(err)
	}
}

// 执行长命令
func RunCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command) // 使用 bash 运行命令
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func RunCommandWithLog(command string) error {
	log.Println("执行命令：", command)
	// 如果是要调用vi的，则需要额外处理，git commit，vi
	if isCallTerminalVim(command) {
		callTerminalVim(command)
		return nil
	}
	// 如果是调用cd命令，使用Chdir进入目录
	if strings.HasPrefix(command, "cd") {
		parts := strings.Fields(command)
		if len(parts) > 1 {
			targetDir := parts[1]
			if strings.Contains(targetDir, "~") {
				homeDir := GetUserHomePath()
				targetDir = strings.ReplaceAll(targetDir, "~", homeDir)
			}
			if err := os.Chdir(targetDir); err != nil {
				fmt.Printf("切换到目录 %s 失败: %v\n", targetDir, err)
			}
		}
		return nil
	}
	cmd := exec.Command("bash", "-c", command) // 使用 bash 运行命令
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
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

// 获取当前网络IP
func GetIpAddress() ([]string, []string) {
	ipv4 := []string{}
	ipv6 := []string{}
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("获取网络接口失败:", err)
		return ipv4, ipv6
	}

	for _, iface := range interfaces {
		// 获取每个接口的地址
		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, address := range addresses {
			// 检查地址类型，过滤掉非 IP 地址
			if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				ipAddr := ipNet.IP.String()
				// 输出非环回地址
				if ipNet.IP.To4() != nil {
					ipv4 = append(ipv4, ipAddr)
				} else {
					ipv6 = append(ipv6, ipAddr)
				}
			}
		}
	}
	return ipv4, ipv6
}

// 复制文本到剪切板
func CopyText2ClipBoard(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		return err
	}
	return nil
}

// 查找最后一个空格前面的内容
// 比如"git push o"，则返回 "git push"
func GetBeforeLastSpace(input string) string {
	// 查找最后一个空格的位置
	lastSpaceIndex := strings.LastIndex(input, " ")
	if lastSpaceIndex == -1 {
		// 如果没有空格，返回原字符串
		return input
	}
	// 返回空格前的部分
	return input[:lastSpaceIndex]
}

// 查询文本中是否包含中文
func ContainsChineseWords(text string) bool {
	// 定义一个匹配中文字符的正则表达式
	reg := regexp.MustCompile("[\u4e00-\u9fa5]")
	return reg.MatchString(text)
}

// 使用对应系统的编辑器，编辑文件
func EditFileWithSystemEditor(filePath string) {
	editorCmd := "code"
	if runtime.GOOS == "linux" {
		editorCmd = "vi"
	}
	RunCommandWithLog(fmt.Sprintf("%s %s", editorCmd, filePath))
}

// TODO：使用对应系统编辑器预览文件
func PreviewFileWithSystemEditor(filePath string) {
	RunCommandWithLog(fmt.Sprintf("cat %s", filePath))
}
