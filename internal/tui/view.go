package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

func (m *Model) View() string {
	return docStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.lists[TaskListsPane].View(),
			m.lists[TasksPane].View(),
		),
	)
}