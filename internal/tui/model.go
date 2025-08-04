package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanicksenn/gtasks/internal/gtasks"
	taskspb "google.golang.org/api/tasks/v1"
)

type taskListsLoadedMsg struct {
	taskLists *taskspb.TaskLists
}

type tasksLoadedMsg struct {
	tasks *taskspb.Tasks
}

type Pane int

const (
	TaskListsPane Pane = iota
	TasksPane
)

type Model struct {
	client  gtasks.Client
	focused Pane
	lists   []list.Model
}

func New() (*Model, error) {
	client, err := gtasks.NewClient(context.Background(), false)
	if err != nil {
		return nil, err
	}

	return &Model{
		client:  client,
		focused: TaskListsPane,
		lists: []list.Model{
			list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
			list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		},
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return func() tea.Msg {
		taskLists, err := m.client.ListTaskLists()
		if err != nil {
			return err
		}
		return taskListsLoadedMsg{taskLists: taskLists}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case taskListsLoadedMsg:
		items := make([]list.Item, len(msg.taskLists.Items))
		for i, taskList := range msg.taskLists.Items {
			items[i] = taskListItem{taskList}
		}
		m.lists[TaskListsPane].SetItems(items)

	case tasksLoadedMsg:
		items := make([]list.Item, len(msg.tasks.Items))
		for i, task := range msg.tasks.Items {
			items[i] = taskItem{task}
		}
		m.lists[TasksPane].SetItems(items)

	case tea.KeyMsg:
		if m.lists[m.focused].FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focused = (m.focused + 1) % 2
		case "enter":
			if m.focused == TaskListsPane {
				selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
				return m, func() tea.Msg {
					tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{TaskListID: selectedTaskList.Id})
					if err != nil {
						return err
					}
					return tasksLoadedMsg{tasks: tasks}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}
