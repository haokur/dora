package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var inputCmd string
var outputCmd string

var generateExecCmd = &cobra.Command{
	Use:   "exe",
	Short: "生成可执行文件；\n例如：dora exe -i \"cd ~ && ls && echo Hello World\" -o get_root_list",
	Run: func(cmd *cobra.Command, args []string) {
		if inputCmd != "" && outputCmd != "" {
			// 指定文件生成的临时目录：用户主目录下的 `~/dora/.cache`
			cacheDir := getTempDirectory()

			// 生成包含执行命令逻辑的代码文件
			goFilePath := filepath.Join(cacheDir, outputCmd+".go")
			generateGoFile(inputCmd, goFilePath)

			// 编译生成的代码为可执行文件
			workDir, _ := os.Getwd()
			executablePath := filepath.Join(workDir, outputCmd)
			compileExecutable(goFilePath, executablePath)

			// 删除临时生成的 .go 文件
			cleanupGoFile(goFilePath)
		} else {
			fmt.Println("Error: Please provide both -i and -o arguments.")
			cmd.Help()
		}
	},
}

// getTempDirectory 获取或创建缓存目录
func getTempDirectory() string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("获取用户目录失败:", err)
		os.Exit(1)
	}

	// 生成缓存目录路径，例如：`~/dora/.cache`
	cacheDir := filepath.Join(homeDir, "dora", ".cache")

	// 检查目录是否存在，若不存在则创建
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(cacheDir, os.ModePerm) // 创建所有必要的目录
		if err != nil {
			fmt.Println("创建缓存目录失败:", err)
			os.Exit(1)
		}
	}

	return cacheDir
}

// generateGoFile 生成包含执行命令逻辑的 Go 源文件
func generateGoFile(command string, filePath string) {
	tpl := `package main

import (
	"log"
	"os"
	"os/exec"
	"fmt"
	"strings"
)

// RunCommandWithLog 执行一条命令并记录日志
func RunCommandWithLog(command string) error {
	log.Println("执行命令：", command)

	// 如果是cd命令，使用Chdir进入目录
	if strings.HasPrefix(command, "cd") {
		parts := strings.Fields(command)
		if len(parts) > 1 {
			targetDir := parts[1]
			if strings.Contains(targetDir, "~") {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					fmt.Println("获取用户目录失败", err)
					return err
				}
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

func main() {
	commands := "{{.Command}}"
	commandList := strings.Split(commands, "&&")
	for _, cmd := range commandList {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" {
			err := RunCommandWithLog(cmd)
			if err != nil {
				log.Fatalf("命令执行失败: %v", err)
			}
		}
	}
}
`

	// 创建并写入 Go 文件
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating Go file:", err)
		return
	}
	defer f.Close()

	// 使用模板将命令写入 Go 文件
	t := template.Must(template.New("goFile").Parse(tpl))
	err = t.Execute(f, struct{ Command string }{Command: command})
	if err != nil {
		fmt.Println("Error writing Go file:", err)
		return
	}

	fmt.Println("Generated Go file:", filePath)
}

// compileExecutable 编译生成的 Go 源文件为可执行文件
func compileExecutable(goFilePath string, executablePath string) {
	cmd := exec.Command("go", "build", "-o", executablePath, goFilePath)

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error generating executable: %s\n", err)
		return
	}

	fmt.Printf("Generated executable: %s\n", executablePath)
}

// cleanupGoFile 删除生成的临时 Go 文件
func cleanupGoFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Printf("Error deleting temporary Go file: %s\n", err)
		return
	}

	fmt.Printf("Deleted temporary Go file: %s\n", filePath)
}

func init() {
	// 定义命令行参数
	generateExecCmd.Flags().StringVarP(&inputCmd, "input", "i", "", "输入要执行的命令，多个命令用&&隔开")
	generateExecCmd.Flags().StringVarP(&outputCmd, "output", "o", "", "输出生成可执行文件名称")
	rootCmd.AddCommand(generateExecCmd)
}
