package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

func (m *Model) View() string {
	if m.state == stateNewTaskList {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			"Create a new task list",
			m.newTaskListInput.View(),
			m.status,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		docStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.lists[TaskListsPane].View(),
				m.lists[TasksPane].View(),
			),
		),
		m.status,
	)
}