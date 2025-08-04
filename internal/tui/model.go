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

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focused = (m.focused + 1) % 2
		}
	}

	return m, nil
}
