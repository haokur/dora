package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// 单项选择的类型
type CommandItem struct {
	Label string
	Value string
}

// 搜索模型
type searchModel struct {
	commands   []CommandItem
	cursor     int
	selected   map[int]struct{}
	filtered   []CommandItem
	searchTerm string
}

// 初始化搜索模型
func initialSearchModel(commands []CommandItem) searchModel {
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
			m.selected = map[int]struct{}{}
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
	s += fmt.Sprintf("输入关键字进行筛选: %s\n\n", m.searchTerm)
	for i, command := range m.filtered {
		cursor := " " // 光标指示符
		if m.cursor == i {
			cursor = ">" // 当前光标所在位置
		}

		checked := " " // 选择状态
		if _, ok := m.selected[i]; ok {
			checked = "√"
		}

		s += fmt.Sprintf("%s [%s] %s（%s）\n", cursor, checked, highlight(command.Value, m.searchTerm), command.Label)
	}

	// 检查是否有已选择的项
	if len(m.selected) > 0 && len(m.filtered) < len(m.selected) {
		s += "\n已选择的项: "
		for k, v := range m.commands {
			if _, ok := m.selected[k]; ok {
				s += v.Value + "；"
			}
		}
	}

	return s
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
func filterCommands(commands []CommandItem, searchTerm string) []CommandItem {
	if searchTerm == "" {
		return commands
	}

	var filtered []CommandItem
	for _, cmd := range commands {
		if isSubsequence(searchTerm, cmd.Value) {
			filtered = append(filtered, cmd)
		}
	}
	return filtered
}

// 高亮匹配字符（按顺序高亮）
func highlight(command, input string) string {
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
			// result.WriteString(fmt.Sprintf("\033[1;31m%s\033[0m", string(command[i]))) // 高亮当前字符
			result.WriteString(fmt.Sprintf("\033[1;31;4m%s\033[0m", string(command[i]))) // 高亮并下划线
			inputIndex++
		} else {
			result.WriteString(string(command[i])) // 正常显示
		}
	}

	return result.String()
}

// 导出的方法
func Search(commands []CommandItem) ([]string, error) {
	p := tea.NewProgram(initialSearchModel(commands))
	result, err := p.Run()
	if err != nil {
		return nil, err
	}

	choices := []string{}
	for k, v := range result.(searchModel).selected {
		fmt.Println(k, v, result.(searchModel).commands[k])
		choices = append(choices, result.(searchModel).commands[k].Value)
	}

	return choices, nil
}
