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
	Desc  string
}

// 搜索模型
type searchModel struct {
	commands   []CommandItem
	cursor     int
	selected   []CommandItem // 直接存储 CommandItem 对象
	filtered   []CommandItem
	searchTerm string
}

// 初始化搜索模型，返回指向 searchModel 的指针
func initialSearchModel(commands []CommandItem) *searchModel {
	return &searchModel{
		commands: commands,
		selected: []CommandItem{},
		filtered: commands,
	}
}

func (m *searchModel) Init() tea.Cmd {
	return nil
}

// 更新模型（键盘输入处理）
func (m *searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// 退出
		case "ctrl+c":
			m.selected = []CommandItem{} // 重置已选择项
			return m, tea.Quit
		// 取消已选
		case "esc":
			m.selected = []CommandItem{} // 重置已选择项
			return m, nil
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
			if isSelected(m.selected, m.filtered[m.cursor]) {
				m.selected = removeSelection(m.selected, m.filtered[m.cursor]) // 取消选择
			} else {
				m.selected = append(m.selected, m.filtered[m.cursor]) // 添加选择
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

// 检查命令是否已被选择
func isSelected(selected []CommandItem, command CommandItem) bool {
	for _, v := range selected {
		if v.Value == command.Value {
			return true
		}
	}
	return false
}

// 移除选择的命令
func removeSelection(selected []CommandItem, command CommandItem) []CommandItem {
	for i, v := range selected {
		if v.Value == command.Value {
			return append(selected[:i], selected[i+1:]...) // 取消选择
		}
	}
	return selected
}

// 渲染界面
// 渲染界面
func (m *searchModel) View() string {
	s := "使用上下键选择，空格选择，回车执行，ctrl+c退出，ESC取消已选\n"
	s += fmt.Sprintf("输入关键字进行筛选: %s\n\n", m.searchTerm)

	for i, command := range m.filtered {
		cursor := " " // 光标指示符
		if m.cursor == i {
			cursor = ">" // 当前光标所在位置
		}

		// 选择状态，前面用顺序数字
		checked := " " // 默认状态
		if isSelected(m.selected, command) {
			// 如果已选择，找到该命令在已选择项中的位置并用数字表示
			for j, selectedCommand := range m.selected {
				if selectedCommand.Value == command.Value {
					checked = fmt.Sprintf("%d", j+1) // 用顺序数字表示
					break
				}
			}
		}

		labelStr := ""
		if command.Label != "" {
			labelStr = fmt.Sprintf("（%s）", command.Label)
		}

		descStr := ""
		if command.Label != "" {
			descStr = darkText(fmt.Sprintf("[%s]", command.Desc))
		}

		s += fmt.Sprintf("%s [%s] %s%s%s\n", cursor, checked, highlight(command.Value, m.searchTerm), labelStr, descStr)
	}

	// 检查是否有已选择的项并展示
	if len(m.selected) > 0 {
		s += "\n已选择的项，将按下面顺序返回:\n"
		for i, cmd := range m.selected {
			s += fmt.Sprintf("%d. %s\n", i+1, cmd.Value) // 显示已选择的命令及其顺序
			if cmd.Desc != "" {
				s += darkText("  - " + strings.ReplaceAll(cmd.Desc, "，", "\n  - "))
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
			result.WriteString(fmt.Sprintf("\033[1;31;4m%s\033[0m", string(command[i]))) // 高亮并下划线
			inputIndex++
		} else {
			result.WriteString(string(command[i])) // 正常显示
		}
	}

	return result.String()
}

func darkText(text string) string {
	return fmt.Sprintf("\033[90m%s\033[0m", text)
}

// 带搜索的多选
func Search(commands []CommandItem) ([]string, error) {
	p := tea.NewProgram(initialSearchModel(commands)) // 传递指向 searchModel 的指针
	result, err := p.Run()
	if err != nil {
		return nil, err
	}

	choices := []string{}
	for _, cmd := range result.(*searchModel).selected {
		choices = append(choices, cmd.Value)
	}

	return choices, nil
}
