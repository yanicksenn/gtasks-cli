package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yanicksenn/gtasks/internal/gtasks"
	"google.golang.org/api/tasks/v1"
)

// The tui package contains the interactive terminal user interface for gtasks.
// It is built using the Bubble Tea framework, which is based on the
// Model-View-Update (MVU) architecture.
//
// The Model is the state of the application.
// The View is a function that renders the state as a string.
// The Update function handles messages (events) and updates the state.
//
// The flow of the application is as follows:
// 1. The user interacts with the UI, which generates a message.
// 2. The message is sent to the Update function.
// 3. The Update function updates the Model.
// 4. The View function re-renders the UI based on the new state.
//
// This cycle repeats until the user quits the application.

// Model is the state of the TUI application.
type Model struct {
	client     gtasks.Client
	err        error
	taskListID string
	taskList   *tasks.TaskList
	tasks      []*tasks.Task
	cursor     int
}

// New creates a new TUI model.
func New(client gtasks.Client, taskListID string) (*Model, error) {
	return &Model{
		client:     client,
		taskListID: taskListID,
	}, nil
}

// Init is the first command that is executed when the application starts.
func (m *Model) Init() tea.Cmd {
	// Fetch the task list and tasks in parallel.
	return tea.Batch(m.fetchTaskList, m.fetchTasks)
}

// Update handles messages and updates the model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle key presses
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
	// Handle fetched task list
	case taskListFetchedMsg:
		m.taskList = msg.taskList
	// Handle fetched tasks
	case tasksFetchedMsg:
		m.tasks = msg.tasks
	// Handle toggled task completion
	case taskCompletionToggledMsg:
		m.tasks[m.cursor] = msg.task
	// Handle errors
	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

// View renders the UI.
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

// renderTask renders a single task.
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

// taskListFetchedMsg is a message that is sent when the task list is fetched.
type taskListFetchedMsg struct {
	taskList *tasks.TaskList
}

// tasksFetchedMsg is a message that is sent when the tasks are fetched.
type tasksFetchedMsg struct {
	tasks []*tasks.Task
}

// taskCompletionToggledMsg is a message that is sent when a task's completion is toggled.
type taskCompletionToggledMsg struct {
	task *tasks.Task
}

// errMsg is a message that is sent when an error occurs.
type errMsg struct {
	err error
}

// fetchTaskList fetches the task list from the client.
func (m *Model) fetchTaskList() tea.Msg {
	taskList, err := m.client.GetTaskList(gtasks.GetTaskListOptions{
		TaskListID: m.taskListID,
	})
	if err != nil {
		return errMsg{err}
	}
	return taskListFetchedMsg{taskList}
}

// fetchTasks fetches the tasks from the client.
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

// toggleTaskCompletion toggles the completion status of the selected task.
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
