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

	mainView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.lists[TaskListsPane].View(),
		tasksView,
	)

	h, _ := docStyle.GetFrameSize()
	detailsViewStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		Padding(1, 2).
		Width(m.width - h).
		Height(m.detailsViewHeight)

	detailsView := ""
	if m.focused == TaskListsPane {
		if selectedItem := m.lists[TaskListsPane].SelectedItem(); selectedItem != nil {
			taskList := selectedItem.(taskListItem)
			detailsView = detailsViewStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					fmt.Sprintf("Title: %s", taskList.Title()),
					fmt.Sprintf("ID: %s", taskList.Id),
				),
			)
		}
	} else if m.focused == TasksPane {
		if selectedItem := m.lists[TasksPane].SelectedItem(); selectedItem != nil {
			task := selectedItem.(taskItem)
			detailsView = detailsViewStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					fmt.Sprintf("Title: %s", task.Title()),
					fmt.Sprintf("Status: %s", task.Status),
					fmt.Sprintf("Notes: %s", task.Description()),
					fmt.Sprintf("Due: %s", task.Due),
				),
			)
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		docStyle.Render(mainView),
		detailsView,
		m.status,
	)
}