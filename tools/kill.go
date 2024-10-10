package tools

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/haokur/dora/cmd"
	ps "github.com/mitchellh/go-ps"
	portNet "github.com/shirou/gopsutil/net"
)

type ProcessItem struct {
	Pid     int
	PPid    int
	Command string
}

// contains 判断切片中是否包含某个元素
func contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// 根据程序名获取Pid信息
func GetPidInfoByPs(processName string) []ProcessItem {
	var pidList []ProcessItem

	// 获取所有进程
	processList, err := ps.Processes()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	lowerCaseProcessName := strings.ToLower(processName)
	// 遍历进程，匹配程序名
	for _, p := range processList {
		execStr := strings.ToLower(p.Executable())
		if execStr == lowerCaseProcessName {
			pidList = append(pidList, ProcessItem{
				Pid:     p.Pid(),
				PPid:    p.PPid(),
				Command: p.Executable(),
			})
		}
	}

	return pidList
}

// 使用ps，从path中匹配
func GetPidInfoByAuxGrep(processName string) []ProcessItem {
	processList := []ProcessItem{}
	switch runtime.GOOS {
	case "windows":
	// TODO:待验证
	case "darwin", "linux":
		// cmd := exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep %s | grep -v grep | awk '{print $2, $3, $11}'", processName))
		cmd := exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep %s | grep -v grep | awk '{print}'", processName))

		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		// 将输出按行分割
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			// 忽略空行
			if strings.TrimSpace(line) == "" {
				continue
			}

			// 使用 strings.Fields 按空白字符切割每一行
			fields := strings.Fields(line)

			if len(fields) >= 11 {
				// 提取所需字段，PID (第2列), PPID (第3列), Command (第11列开始)
				pid, _ := strconv.Atoi(fields[1])
				ppid, _ := strconv.Atoi(fields[2])
				command := strings.Join(fields[10:], " ") // 命令可能包含空格，所以从第11列开始

				if !strings.Contains(command, fmt.Sprintf("kill %s", processName)) {
					processList = append(processList, ProcessItem{
						Pid:     pid,
						PPid:    ppid,
						Command: command,
					})
				}
			}
		}
	default:
		fmt.Println("Unsupported operating system")
	}
	return processList
}

// 根据IP获取pid信息
func GetPidInfoByPort(port int) []ProcessItem {
	pidList := []ProcessItem{}
	// 获取所有网络连接
	connections, _ := portNet.Connections("inet")

	_port := uint32(port)
	// 遍历连接，匹配端口号
	for _, conn := range connections {
		if conn.Laddr.Port == _port {
			pidList = append(pidList, ProcessItem{
				Pid:     int(conn.Pid),
				PPid:    1,
				Command: fmt.Sprintf(":%d", port),
			})
		}
	}

	return pidList
}

// 命令行交互-用户选择要kill的匹配的进程
func selectPid2Kill(pidList *[]ProcessItem, processName string) []int {
	killPidList := []int{}

	selectOptions := []string{}
	for _, v := range *pidList {
		// optionItem := fmt.Sprintf("COMMAND：%s，PID：%s，NAME：%s", v["COMMAND"], v["PID"], v["NAME"])
		optionItem := fmt.Sprintf("PID：%d，PPID：%d，COMMAND：%s", v.Pid, v.PPid, v.Command)
		selectOptions = append(selectOptions, optionItem)
	}
	// 拼接选择项
	_, allChoiceIndex, err := cmd.Check(fmt.Sprintf("选择对应【%s】要kill的进程", processName), &selectOptions, false)
	if err != nil {
		fmt.Println("选择要kill的进程出错://", err)
	}

	for index, v := range *pidList {
		if contains(allChoiceIndex, index) {
			killPidList = append(killPidList, v.Pid)
		}
	}

	return killPidList
}

// 最终都是用pid来kill进程
func killProcessByPid(pidArr []int, processName string) {
	for _, pid := range pidArr {
		process, err := os.FindProcess(pid)
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
func KillProcess(args *[]string, silence bool) {
	for _, processItem := range *args {
		pidInfoList := []ProcessItem{}
		port, err := strconv.Atoi(processItem)
		if err != nil {
			pidInfoList = append(pidInfoList, GetPidInfoByPs(processItem)...)

			auxPidList := GetPidInfoByAuxGrep(processItem)
			seen := make(map[string]bool)
			var uniqueProcesses []ProcessItem
			for _, process := range auxPidList {
				if !seen[string(process.Pid)] {
					uniqueProcesses = append(uniqueProcesses, process)
					seen[string(process.Pid)] = true
				}
			}
			pidInfoList = append(pidInfoList, uniqueProcesses...)
		} else {
			pidInfoList = append(pidInfoList, GetPidInfoByPort(port)...)
		}

		if silence {
			willKillPidList := []int{}
			for _, v := range pidInfoList {
				willKillPidList = append(willKillPidList, v.Pid)
			}
			killProcessByPid(willKillPidList, processItem)
		} else {
			selectPidList := selectPid2Kill(&pidInfoList, processItem)

			if len(selectPidList) > 0 {
				killProcessByPid(selectPidList, processItem)
			}
		}
	}
}
