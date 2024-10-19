package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/haokur/dora/cmd"
)

// 获取未提交的文件列表
func getUncommittedFiles(sourceDir string) ([]string, error) {
	os.Chdir(sourceDir)
	// git ls-files -m && git ls-files --others --exclude-standard
	// cmd := exec.Command("git", "status", "--porcelain")
	// out, err := cmd.Output()
	out, err := RunCommand("git ls-files -m && git ls-files --others --exclude-standard")
	if err != nil {
		return nil, fmt.Errorf("failed to get uncommitted files: %w", err)
	}

	// var files []string
	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	// for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
	// 	files = append(files, line)
	// 	// if len(line) > 3 {
	// 	// 	files = append(files, strings.TrimSpace(line[3:])) // 截取文件路径部分
	// 	// }
	// }
	return files, nil
}

// 复制文件到目标路径，保持层级结构
func copyFile(src, dest string) error {
	// 创建目标文件夹
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// 创建目标文件
	destinationFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	// 复制内容
	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// 备份未提交的文件
func backupUncommittedFiles(sourceDir string, backupDir string) error {
	files, err := getUncommittedFiles(sourceDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		source := filepath.Join(sourceDir, file)
		dest := filepath.Join(backupDir, file)
		fmt.Printf("Backing up: %s -> %s\n", source, dest)
		if err := copyFile(source, dest); err != nil {
			return fmt.Errorf("failed to backup file %s: %w", source, err)
		}
	}

	return nil
}

// 将当前sourceDir目录下的更改的文件，以时间戳为文件夹名备份到targetDir
func BackupUnCommitFiles(sourceDir string, targetDir string) (string, error) {
	// 获取当前时间戳
	timestamp := time.Now().Format("2006_01_02_150405")

	// 创建备份文件夹
	// backupDir := filepath.Join(targetDir, timestamp)
	backupDir := targetDir + "_" + timestamp
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		fmt.Println("Error creating backup directory:", err)
		return "", err
	}

	// 备份未提交的文件
	if err := backupUncommittedFiles(sourceDir, backupDir); err != nil {
		fmt.Println("Error backing up files:", err)
		return "", err
	}

	fmt.Printf("Backup completed successfully.\nBackup dir is %s\n", backupDir)
	return backupDir, nil
}

// 将当前backupDir以时间戳为文件夹下所有文件还原到git项目目录下
// 1.找到匹配的备份目录
// 2.以时间戳按时间倒序，最近的备份显示在最前面，单选
// 3.用户选择一个备份目录，点击确认
// 4.展示选择备份目录下所有文件，且显示更改时间，文件大小，用户选择要还原的文件
// 5.将用户选择的文件，还原到git项目目录
func RecoverBackupFiles(backupDir string, gitProjectDir string) {
	backupItemList, err := os.ReadDir(backupDir)
	if err != nil {
		fmt.Println("读取备份目录目录时出错", err)
	}

	// 过滤出文件夹
	var dirs []string
	for _, entry := range backupItemList {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	// 按文件夹名称倒序排序
	// sort.Slice(dirs, func(i, j int) bool {
	// 	return dirs[i] > dirs[j] // 倒序排列
	// })
	// 按文件中的时间戳倒序排序
	dirs = SortSliceByInlineDate(dirs, "2006_01_02_150405", false)

	userSelectBackupDir, err := cmd.Radio("请选择一个文件夹进行还原", &dirs)
	if err != nil {
		fmt.Println("用户选择目录出错", err)
	}
	if userSelectBackupDir == "" {
		return
	}
	backupDir2Recover := filepath.Join(backupDir, userSelectBackupDir)
	allFilePaths, err := ReadFilesRecursively(backupDir2Recover)

	if err != nil {
		fmt.Println("递归读取报错", err)
	}
	// 提示用户选择要还原的文件
	userSelectFiles, _, err := cmd.Check("请选择要还原的文件", &allFilePaths, false)
	if err != nil {
		fmt.Println("用户选择文件出错", err)
	}
	// 将用户选择的还原到git项目目录
	for _, recoverFilePath := range userSelectFiles {
		source := filepath.Join(backupDir2Recover, recoverFilePath)
		dest := filepath.Join(gitProjectDir, recoverFilePath)
		fmt.Printf("Backing up: %s -> %s\n", source, dest)
		copyFile(source, dest)
	}
	fmt.Println("recover successfully!")
}
