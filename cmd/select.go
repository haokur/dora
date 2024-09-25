package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type selectModel struct {
	label       string       // 标题标签
	emptyEnable bool         // 是否允许为空
	choices     []string     // 可供选择的选项
	cursor      int          // 当前光标位置
	checked     map[int]bool // 保存已选择的选项
	done        bool         // 用户是否完成选择
	allSelected bool         // 标记是否全选
	isCanceled  bool         // 是否取消
}

func initialSelectModel(label string, options *[]string, emptyEnable bool) selectModel {
	return selectModel{
		label:       label,
		choices:     *options,
		checked:     make(map[int]bool),
		emptyEnable: emptyEnable,
	}
}

func (m selectModel) Init() tea.Cmd {
	// 初始化时不需要做什么操作
	return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// 处理键盘事件
	case tea.KeyMsg:
		switch msg.String() {

		// 上下箭头用于移动光标
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// 空格键用于选择或取消选择
		case " ":
			_, ok := m.checked[m.cursor]
			if ok {
				delete(m.checked, m.cursor) // 如果已经选中，则取消选中
			} else {
				m.checked[m.cursor] = true // 否则标记为选中
			}

		// Enter 键用于提交选择
		case "enter":
			canSubmit := m.emptyEnable || (!m.emptyEnable && len(m.checked) != 0)
			if canSubmit {
				m.done = true
				return m, tea.Quit
			}

		// 全选或取消全选 (按键 'a')
		case "a":
			if m.allSelected {
				// 如果已经全选，执行取消全选
				m.checked = make(map[int]bool)
				m.allSelected = false
			} else {
				// 全选
				for i := range m.choices {
					m.checked[i] = true
				}
				m.allSelected = true
			}

		// 退出
		case "ctrl+c", "q":
			m.isCanceled = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// 返回用户选择的结果
func (m selectModel) SelectedChoices() string {
	selectedChoice := []string{}
	for i, choice := range m.choices {
		if m.checked[i] {
			selectedChoice = append(selectedChoice, choice)
		}
	}
	// 这里最后要加\n，不然可能显示不出来
	return fmt.Sprintf("%s: %s\n", m.label, strings.Join(selectedChoice, ","))
}

func (m selectModel) View() string {
	if m.isCanceled {
		return fmt.Sprintf("%s: %s\n", m.label, "操作已取消")
	}

	if m.done {
		// 完成选择时显示结果
		return m.SelectedChoices()
	}

	// 构建选择列表的界面
	// s := "请选择选项 (空格选择，a全选/取消全选，Enter提交)：\n\n"
	s := fmt.Sprintf("%s (空格选择，a全选/取消全选，Enter提交)：\n\n", m.label)
	// s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)

	for i, choice := range m.choices {

		// 显示光标
		cursor := " " // 未选中项前面显示空格
		if m.cursor == i {
			cursor = ">" // 光标位置的项前显示 >
		}

		// 显示已选中的选项
		checked := " " // 默认未选中
		if m.checked[i] {
			checked = "√" // 已选中的项前显示 x
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\n按 q 退出\n"
	return s
}

// 多选
func Check(label string, options *[]string, emptyEnable bool) ([]string, []int, error) {
	p := tea.NewProgram(initialSelectModel(label, options, emptyEnable))
	allChoice := []string{}
	allChoiceIndex := []int{}
	result, err := p.Run()
	if err != nil {
		return allChoice, []int{}, err
	}
	for i, choice := range result.(selectModel).choices {
		if result.(selectModel).checked[i] {
			allChoice = append(allChoice, choice)
			allChoiceIndex = append(allChoiceIndex, i)
		}
	}
	return allChoice, allChoiceIndex, nil
}
