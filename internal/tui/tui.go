package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yanicksenn/gtasks/internal/gtasks"
	"google.golang.org/api/tasks/v1"
)

type Model struct {
	client     gtasks.Client
	err        error
	taskListID string
	taskList   *tasks.TaskList
	tasks      []*tasks.Task
	cursor     int
}

func New(client gtasks.Client, taskListID string) (*Model, error) {
	return &Model{
		client:     client,
		taskListID: taskListID,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.fetchTaskList, m.fetchTasks)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case " ":
			return m, m.toggleTaskCompletion
		}
	case taskListFetchedMsg:
		m.taskList = msg.taskList
	case tasksFetchedMsg:
		m.tasks = msg.tasks
	case taskCompletionToggledMsg:
		m.tasks[m.cursor] = msg.task
	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

func (m *Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %s", m.err)
	}

	if m.taskList == nil {
		return "Loading..."
	}

	s := fmt.Sprintf("%s\n\n", m.taskList.Title)
	for i, task := range m.tasks {
		s += m.renderTask(i, task)
	}
	s += "\nPress q to quit."
	return s
}

func (m *Model) renderTask(i int, task *tasks.Task) string {
	checkbox := "[ ]"
	if task.Status == "completed" {
		checkbox = "[x]"
	}

	style := lipgloss.NewStyle()
	if m.cursor == i {
		style = style.Foreground(lipgloss.Color("205"))
	}

	return fmt.Sprintf("%s\n", style.Render(fmt.Sprintf("%s %s", checkbox, task.Title)))
}

type taskListFetchedMsg struct {
	taskList *tasks.TaskList
}

type tasksFetchedMsg struct {
	tasks []*tasks.Task
}

type taskCompletionToggledMsg struct {
	task *tasks.Task
}

type errMsg struct {
	err error
}

func (m *Model) fetchTaskList() tea.Msg {
	taskList, err := m.client.GetTaskList(gtasks.GetTaskListOptions{
		TaskListID: m.taskListID,
	})
	if err != nil {
		return errMsg{err}
	}
	return taskListFetchedMsg{taskList}
}

func (m *Model) fetchTasks() tea.Msg {
	tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{
		TaskListID:    m.taskListID,
		ShowCompleted: true,
		SortBy:        "alphabetical",
	})
	if err != nil {
		return errMsg{err}
	}
	return tasksFetchedMsg{tasks.Items}
}

func (m *Model) toggleTaskCompletion() tea.Msg {
	task := m.tasks[m.cursor]

	var updatedTask *tasks.Task
	var err error

	if task.Status == "completed" {
		updatedTask, err = m.client.UncompleteTask(gtasks.UncompleteTaskOptions{
			TaskListID: m.taskListID,
			TaskID:     task.Id,
		})
	} else {
		updatedTask, err = m.client.CompleteTask(gtasks.CompleteTaskOptions{
			TaskListID: m.taskListID,
			TaskID:     task.Id,
		})
	}

	if err != nil {
		return errMsg{err}
	}

	return taskCompletionToggledMsg{updatedTask}
}
