package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	textInput    textinput.Model
	confirmed    bool
	label        string
	defaultValue string
	value        string
	isCanceled   bool // 是否取消
}

func initialInputModel(label string, defaultValue string) inputModel {
	ti := textinput.New()
	ti.Placeholder = defaultValue // 提示符号
	ti.Focus()                    // 聚焦输入
	ti.Width = 10                 // 设置宽度

	return inputModel{
		textInput:    ti,
		label:        label,
		defaultValue: defaultValue,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink // 启动光标闪烁
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// 回车键确认输入
		case "enter":
			inputValue := m.textInput.Value()
			if inputValue == "" {
				m.textInput.SetValue(m.defaultValue)
				inputValue = m.defaultValue
			}
			m.value = inputValue
			m.confirmed = true
			return m, tea.Quit

		// 退出
		case "ctrl+c":
			m.isCanceled = true
			return m, tea.Quit

		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg) // 更新输入
			return m, cmd
		}
	}
	return m, nil
}

func (m inputModel) View() string {
	if m.isCanceled {
		return fmt.Sprintf("%s: %s\n", m.label, "操作已取消")
	}
	if m.confirmed {
		return fmt.Sprintf("%s: %s\n", m.label, m.textInput.Value())
	}
	return fmt.Sprintf("%s %s", m.label, m.textInput.View()) // 单行显示提示和光标输入
}

func Input(label string, defaultValue string) (string, error) {
	p := tea.NewProgram(initialInputModel(label, defaultValue))
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	return result.(inputModel).value, nil
}
