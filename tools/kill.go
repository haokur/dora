package tools

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/haokur/dora/cmd"
)

// contains 判断切片中是否包含某个元素
func contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// 执行lsof 获取
func getPidListByRunLsof(lsofCommand *exec.Cmd) []map[string]string {
	lsofOutput, err := lsofCommand.Output()
	if err != nil {
		fmt.Println("[error]", lsofCommand, "Error executing lsof command,maybe cannot found port or process: ", err)
		return []map[string]string{}
	}

	lsofOutputStr := string(lsofOutput)
	lsofOutputRows := strings.Split(lsofOutputStr, "\n")
	labelRow := lsofOutputRows[0]
	labelArr := strings.Fields(labelRow)

	valueRows := lsofOutputRows[1:]

	lsofProcessList := []map[string]string{}
	pidOnlyList := []string{}
	for _, v := range valueRows {
		valueRowArr := strings.Fields(v)
		if len(valueRowArr) != 0 {
			rowMap := map[string]string{}
			for index, rowMapKey := range labelArr {
				if len(valueRowArr) > index {
					rowMap[rowMapKey] = valueRowArr[index]
				}
			}
			// 过滤掉重复的pid
			currentPid := rowMap["PID"]
			if contains(pidOnlyList, currentPid) {
				continue
			}
			pidOnlyList = append(pidOnlyList, currentPid)
			lsofProcessList = append(lsofProcessList, rowMap)
		}
	}

	return lsofProcessList
}

// 执行ps aux获取
func getPidListByPsAux(processName string) []map[string]string {
	processList := []map[string]string{}
	switch runtime.GOOS {
	case "windows":
		// TODO:待验证
		output, err := RunCommand("netstat -ano")
		if err != nil {
			fmt.Println(err)
		}

		for _, v := range strings.Split(string(output), "\n") {
			if v != "" {
				processList = append(processList, map[string]string{
					"PID":     v,
					"COMMAND": processName,
					"NAME":    processName,
				})
			}
		}
	case "darwin", "linux":
		cmd := exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep %s | grep -v grep | awk '{print $2}'", processName))

		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		for _, v := range strings.Split(string(output), "\n") {
			if v != "" {
				processList = append(processList, map[string]string{
					"PID":     v,
					"COMMAND": processName,
					"NAME":    processName,
				})
			}
		}
	default:
		fmt.Println("Unsupported operating system")
	}
	return processList
}

// 命令行交互-用户选择要kill的匹配的进程
func selectPid2Kill(pidList *[]map[string]string, processName string) []string {
	killPidList := []string{}

	selectOptions := []string{}
	for _, v := range *pidList {
		optionItem := fmt.Sprintf("COMMAND：%s，PID：%s，NAME：%s", v["COMMAND"], v["PID"], v["NAME"])
		selectOptions = append(selectOptions, optionItem)
	}
	// 拼接选择项
	_, allChoiceIndex, err := cmd.Check(fmt.Sprintf("选择对应【%s】要kill的进程", processName), &selectOptions, false)
	if err != nil {
		fmt.Println("选择要kill的进程出错://", err)
	}

	for index, v := range *pidList {
		if contains(allChoiceIndex, index) {
			killPidList = append(killPidList, v["PID"])
		}
	}

	return killPidList
}

// 最终都是用pid来kill进程
func killProcessByPid(pidArr []string, processName string) {
	for _, pid := range pidArr {
		pidInt, err := strconv.Atoi(pid)
		if err != nil {
			fmt.Println("pid to int fail", err)
			continue
		}
		process, err := os.FindProcess(pidInt)
		if err != nil {
			fmt.Println("FindProcess error:", err)
		}
		process.Kill()
		process.Wait()
	}
	fmt.Printf("%s kill successfully\n", processName)
}

// 传入端口和应用程序的字符串数组
// 如：KillProcess(&[]string{"5173", "obsidian"})
func KillProcess(args *[]string) {
	for _, processItem := range *args {
		_, err := strconv.Atoi(processItem)
		lsofParams := []string{}
		pidList := []map[string]string{}
		if err != nil {
			// 如果是应用，使用ps aux查找一遍
			pidList = append(pidList, getPidListByPsAux(processItem)...)
			lsofParams = []string{"-c", processItem}
		} else {
			lsofParams = []string{"-i", fmt.Sprintf(":%s", processItem)}
		}

		lsofCmd := exec.Command("lsof", lsofParams...)
		pidList = append(pidList, getPidListByRunLsof(lsofCmd)...)
		if len(pidList) == 0 {
			continue
		}
		selectPidList := selectPid2Kill(&pidList, processItem)
		if len(selectPidList) > 0 {
			killProcessByPid(selectPidList, processItem)
		}
	}
}
