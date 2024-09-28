package cmd

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// 搜索模型
type searchModel struct {
	commands   []string
	cursor     int
	selected   map[int]struct{}
	filtered   []string
	searchTerm string
}

// 初始化搜索模型
func initialSearchModel(commands []string) searchModel {
	return searchModel{
		commands: commands,
		selected: make(map[int]struct{}),
		filtered: commands,
	}
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

// 更新模型（键盘输入处理）
func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case " ":
			// 选择/取消选择命令
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			return m, tea.Quit
		case "backspace":
			if len(m.searchTerm) > 0 {
				m.searchTerm = m.searchTerm[:len(m.searchTerm)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.searchTerm += msg.String()
			}
		}

		// 根据搜索词过滤命令
		m.filtered = filterCommands(m.commands, m.searchTerm)
	}

	return m, nil
}

// 渲染界面
func (m searchModel) View() string {
	s := "使用上下键选择命令，按空格选择，回车执行，ESC 退出\n"
	s += fmt.Sprintf("输入关键字进行筛选: %s\n", m.searchTerm)
	for i, command := range m.filtered {
		cursor := " " // cursor 指示符
		if m.cursor == i {
			cursor = ">" // 当前光标所在位置
		}

		checked := " " // 选择状态
		if _, ok := m.selected[i]; ok {
			checked = "√"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, command)
	}

	return s
}

// 执行选中的命令
func (m searchModel) executeCommands() {
	for i := range m.selected {
		cmd := m.filtered[i]
		fmt.Printf("执行命令: %s\n", cmd)

		// 使用 exec.Command 执行系统命令
		command := exec.Command("bash", "-c", cmd)
		output, err := command.CombinedOutput()
		if err != nil {
			fmt.Printf("命令执行出错: %v\n", err)
		}
		fmt.Println(string(output))
	}
}

// isSubsequence 检查短字符串是否为长字符串的子序列
func isSubsequence(short, long string) bool {
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

// 根据搜索词过滤命令
func filterCommands(commands []string, searchTerm string) []string {
	if searchTerm == "" {
		return commands
	}

	var filtered []string
	for _, cmd := range commands {
		if isSubsequence(searchTerm, cmd) {
			filtered = append(filtered, cmd)
		}
	}
	return filtered
}

// 导出的方法
func Search(commands []string) ([]string, error) {
	p := tea.NewProgram(initialSearchModel(commands))
	result, err := p.Run()
	if err != nil {
		return nil, err
	}

	choices := []string{}
	for k, v := range result.(searchModel).selected {
		fmt.Println(k, v, result.(searchModel).commands[k])
		choices = append(choices, result.(searchModel).commands[k])
	}

	return choices, nil
}
