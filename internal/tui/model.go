package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Pane int

const (
	TaskListsPane Pane = iota
	TasksPane
)

type Model struct {
	focused Pane
	lists   []list.Model
}

func New() *Model {
	return &Model{
		focused: TaskListsPane,
		lists: []list.Model{
			list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
			list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
