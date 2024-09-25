package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type confirmModel struct {
	textInput     textinput.Model
	label         string
	answer        string
	defaultAnswer string
	confirmed     bool
	isCanceled    bool // 是否取消
}

func initialConfirmModel(label string, defaultAnswer bool) confirmModel {
	_defaultAnswer := ""
	_labelStr := label

	// 初始化 textinput 组件
	ti := textinput.New()
	if defaultAnswer {
		ti.Placeholder = "Y"
		_defaultAnswer = "Y"
		_labelStr += "（Y/n）"
	} else {
		ti.Placeholder = "N"
		_defaultAnswer = "N"
		_labelStr += "（y/N）"
	}
	ti.Focus()       // 聚焦输入
	ti.CharLimit = 1 // 限制输入长度为1个字符
	ti.Width = 10    // 设置宽度

	return confirmModel{
		textInput:     ti,
		label:         _labelStr,
		confirmed:     false,
		defaultAnswer: _defaultAnswer,
	}
}

func (m confirmModel) Init() tea.Cmd {
	return textinput.Blink // 启动光标闪烁
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// 回车键确认输入
		case "enter":
			inputValue := m.textInput.Value()
			if inputValue == "" {
				inputValue = m.defaultAnswer
				m.textInput.SetValue(inputValue)
			}
			inputValue2Upper := strings.ToUpper(inputValue)
			if inputValue2Upper == "Y" || inputValue2Upper == "N" {
				m.confirmed = true
				m.answer = inputValue
				return m, tea.Quit
			}

		// 退出
		case "ctrl+c", "q":
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

func (m confirmModel) View() string {
	if m.isCanceled {
		return fmt.Sprintf("%s: %s\n", m.label, "操作已取消")
	}

	if m.confirmed {
		answerStr := m.answer
		if (m.defaultAnswer == "Y" && answerStr == "y") || (m.defaultAnswer == "N" && answerStr == "n") {
			answerStr = strings.ToUpper(answerStr)
		}
		return fmt.Sprintf("%s: %s\n", m.label, answerStr)
	}
	return fmt.Sprintf("%s%s\n", m.label, m.textInput.View())
}

func Confirm(label string, defaultChoice bool) (bool, error) {
	p := tea.NewProgram(initialConfirmModel(label, defaultChoice))
	result, err := p.Run()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return strings.ToUpper(result.(confirmModel).answer) == "Y", nil
}
