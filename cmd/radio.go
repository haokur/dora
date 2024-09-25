package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// 模型定义
type radioModel struct {
	cursor     int      // 当前选中的索引
	label      string   // 标题标签
	choices    []string // 可选项
	selected   string   // 最终选择的项
	isCanceled bool     // 是否取消
}

func (m radioModel) Init() tea.Cmd {
	return nil
}

// 初始化
func initialRadioModel(label string, options *[]string) radioModel {
	return radioModel{
		label:   label,
		choices: *options,
	}
}

// 更新函数处理输入
func (m radioModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// 处理键盘输入
	case tea.KeyMsg:
		switch msg.String() {

		// 上下键导航
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// 回车确认
		case "enter":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit // 退出程序

		// 退出
		case "ctrl+c", "q":
			m.isCanceled = true
			return m, tea.Quit
		}

	}
	return m, nil
}

// 渲染函数
func (m radioModel) View() string {
	if m.isCanceled {
		return fmt.Sprintf("%s: %s\n", m.label, "操作已取消")
	}
	s := m.label + "（使用上下键导航，按回车确认选择）：\n\n"

	// 列出所有选项
	for i, choice := range m.choices {
		cursor := " " // 默认无光标
		if m.cursor == i {
			cursor = ">" // 当前选中项前加光标
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	if m.selected != "" {
		s = fmt.Sprintf("%s: %s\n", m.label, m.selected)
	} else {
		s += "\n按 q 退出\n"
	}

	return s
}

// 单选
func Radio(label string, options *[]string) (string, error) {
	p := tea.NewProgram(initialRadioModel(label, options))
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	return result.(radioModel).selected, nil
}
