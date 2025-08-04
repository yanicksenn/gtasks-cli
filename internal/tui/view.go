package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle       = lipgloss.NewStyle().Margin(1, 2)
	emptyListStyle = lipgloss.NewStyle().Margin(4, 2, 1, 2).Foreground(lipgloss.Color("240"))
)

func (m *Model) View() string {
	if m.state == stateTaskView {
		return lipgloss.Place(
			100, 100,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), true).
				Padding(1, 2).
				Render(
					lipgloss.JoinVertical(
						lipgloss.Left,
						fmt.Sprintf("Title: %s", m.selectedTask.Title()),
						fmt.Sprintf("Status: %s", m.selectedTask.Status),
						fmt.Sprintf("Notes: %s", m.selectedTask.Description()),
						fmt.Sprintf("Due: %s", m.selectedTask.Due),
					),
				),
		)
	}

	if m.state == stateNewTaskList {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			"Create a new task list",
			m.newTaskListInput.View(),
			m.status,
		)
	}

	if m.state == stateDeleteTaskList || m.state == stateDeleteTask {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			"Are you sure you want to delete this?",
			m.status,
		)
	}

	tasksView := m.lists[TasksPane].View()
	if len(m.lists[TasksPane].Items()) == 0 {
		emptyText := "Empty task list"
		if m.lists[TasksPane].Title == "Loading..." {
			emptyText = "Loading..."
		}
		tasksView = emptyListStyle.Render(emptyText)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		docStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.lists[TaskListsPane].View(),
				tasksView,
			),
		),
		m.status,
	)
}