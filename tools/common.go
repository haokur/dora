package tools

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

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

// 获取配置文件路径
func GetDoraConfigPath() string {
	userHomeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(userHomeDir, "dora/.config.json")
	return configPath
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
		if entry.IsDir() && !strings.Contains(entry.Name(), "node_modules") {
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

// 浅层读取文件夹下所有文件
type IFileItem struct {
	Path         string
	Name         string
	Size         int64
	LastModified time.Time
}

// 读取第一层
func ReadFilesShallowly(dirPath string) ([]IFileItem, error) {
	// 读取目录下的所有条目（文件和子目录）
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	result := []IFileItem{}

	// 遍历条目
	for _, entry := range entries {
		// 检查是否是文件
		if !entry.IsDir() {
			// 获取完整的文件路径
			fullPath := filepath.Join(dirPath, entry.Name())

			// 使用 os.Stat 获取文件的详细信息
			info, err := os.Stat(fullPath)
			if err != nil {
				log.Println(err)
				continue
			}

			// 获取文件大小和最后修改时间
			size := info.Size()       // 文件大小 (字节)
			modTime := info.ModTime() // 文件的最后修改时间

			result = append(result, IFileItem{
				Path:         fullPath,
				Name:         entry.Name(),
				Size:         size,
				LastModified: modTime,
			})
		}
	}

	return result, nil
}

// 排序文件的函数
func SortFiles(files []IFileItem, sortBy string) {
	switch sortBy {
	case "name":
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name > files[j].Name
		})
	case "size":
		sort.Slice(files, func(i, j int) bool {
			return files[i].Size > files[j].Size
		})
	case "modtime":
		sort.Slice(files, func(i, j int) bool {
			return files[i].LastModified.After(files[j].LastModified)
		})
	default:
		fmt.Println("Unknown sort method, using name by default")
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name < files[j].Name
		})
	}
}

// 反转文件顺序的函数
func ReverseSlice[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
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

// 检查切片中是否包含指定的元素
func SliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
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

// OpenFolderAndSelectFile 打开文件夹并高亮显示指定文件
func OpenFolderAndSelectFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "windows":
		// Windows: 使用 `explorer` 打开并选择文件
		cmd := exec.Command("explorer", "/select,", absPath)
		return cmd.Run()
	case "darwin":
		// macOS: 使用 `open` 命令打开文件夹并高亮文件
		cmd := exec.Command("open", "-R", absPath)
		return cmd.Run()
	case "linux":
		// 对于大部分Linux发行版，可以使用xdg-open打开文件夹，但没有选择文件的功能
		cmd := exec.Command("xdg-open", filepath.Dir(absPath))
		return cmd.Run()
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

// sortedItem 结构体存储字符串和时间戳
type sortedItem struct {
	value     string
	timestamp time.Time
}

// generateRegexFromDateFormat 根据日期格式生成正则表达式
func generateRegexFromDateFormat(dateFormat string) string {
	// 规则映射，将日期格式中的符号映射为对应的正则表达式
	rules := map[string]string{
		"2006": `\d{4}`, // 年
		"01":   `\d{2}`, // 月
		"02":   `\d{2}`, // 日
		"15":   `\d{2}`, // 时
		"04":   `\d{2}`, // 分
		"05":   `\d{2}`, // 秒
	}

	// 用于构建正则表达式的字符串
	regexPattern := dateFormat

	// 替换日期格式中的符号为对应的正则表达式
	for key, pattern := range rules {
		regexPattern = strings.ReplaceAll(regexPattern, key, pattern)
	}

	return regexPattern
}

// SortSliceByInlineDate 按字符串中的日期排序
// 参数：
// - slice: 需要排序的字符串slice
// - dateFormat: 日期时间的格式（例如 "2006_01_02_150405"）
// - ascending: 如果为 true 则正序排序，否则倒序排序
func SortSliceByInlineDate(slice []string, dateFormat string, ascending bool) []string {
	// 生成匹配日期的正则表达式
	regexPattern := generateRegexFromDateFormat(dateFormat)

	// 定义正则表达式
	timePattern := regexp.MustCompile(regexPattern)

	// 存储字符串和时间的映射
	var items []sortedItem

	// 遍历字符串slice，提取时间戳并解析
	for _, str := range slice {
		matches := timePattern.FindString(str)
		if matches != "" {
			// 将匹配的时间戳转换为标准时间格式
			timestamp, err := time.Parse(dateFormat, matches)
			if err == nil {
				items = append(items, sortedItem{value: str, timestamp: timestamp})
			}
		}
	}

	// 按时间排序
	sort.Slice(items, func(i, j int) bool {
		if ascending {
			return items[i].timestamp.Before(items[j].timestamp)
		}
		return items[i].timestamp.After(items[j].timestamp)
	})

	// 提取排序后的结果
	var sortedSlice []string
	for _, item := range items {
		sortedSlice = append(sortedSlice, item.value)
	}

	return sortedSlice
}

func FormatSize(size int64) string {
	// 定义单位
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	// 根据大小选择合适的单位进行转换
	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// 复制文件，源路径复制到目标路径
func CopyFile(srcPath, dstPath string) error {
	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// 确保文件写入完成
	err = dstFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}
