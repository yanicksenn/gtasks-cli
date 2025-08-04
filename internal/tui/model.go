package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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

type taskListCreatedMsg struct {
	taskList *taskspb.TaskList
}

type taskDeletedMsg struct{}

type errorMsg struct {
	err error
}

func (e errorMsg) Error() string {
	return e.err.Error()
}

type state int

const (
	stateDefault state = iota
	stateNewTaskList
	stateDeleteTaskList
	stateDeleteTask
)

type Pane int

const (
	TaskListsPane Pane = iota
	TasksPane
)

type Model struct {
	client         gtasks.Client
	focused        Pane
	lists          []list.Model
	status         string
	state          state
	newTaskListInput textinput.Model
}

func New() (*Model, error) {
	client, err := gtasks.NewClient(context.Background(), false)
	if err != nil {
		return nil, err
	}

	newTaskListInput := textinput.New()
	newTaskListInput.Placeholder = "New Task List"
	newTaskListInput.Focus()

	m := &Model{
		client:         client,
		focused:        TaskListsPane,
		lists: []list.Model{
			list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
			list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		},
		state:          stateDefault,
		newTaskListInput: newTaskListInput,
	}
	m.SetStatus("Ready")
	return m, nil
}

func (m *Model) SetStatus(status string) {
	m.status = status
}

func (m *Model) Init() tea.Cmd {
	return func() tea.Msg {
		taskLists, err := m.client.ListTaskLists()
		if err != nil {
			return errorMsg{err}
		}
		return taskListsLoadedMsg{taskLists: taskLists}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.lists[TaskListsPane].SetSize(msg.Width/2-h, msg.Height-v)
		m.lists[TasksPane].SetSize(msg.Width/2-h, msg.Height-v)

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

	case taskListCreatedMsg:
		m.newTaskListInput.Reset()
		m.state = stateDefault
		m.SetStatus("Task list created")
		return m, m.Init()

	case taskDeletedMsg:
		m.SetStatus("Task deleted")
		selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
		return m, func() tea.Msg {
			tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{TaskListID: selectedTaskList.Id})
			if err != nil {
				return errorMsg{err}
			}
			return tasksLoadedMsg{tasks: tasks}
		}

	case errorMsg:
		m.SetStatus(msg.Error())
		return m, nil

	case tea.KeyMsg:
		if m.state == stateNewTaskList {
			switch keypress := msg.String(); keypress {
			case "enter":
				title := m.newTaskListInput.Value()
				m.state = stateDefault
				m.SetStatus("Creating task list...")
				return m, func() tea.Msg {
					taskList, err := m.client.CreateTaskList(gtasks.CreateTaskListOptions{Title: title})
					if err != nil {
						return errorMsg{err}
					}
					return taskListCreatedMsg{taskList: taskList}
				}
			case "esc":
				m.state = stateDefault
				m.SetStatus("Ready")
			}
		}

		if m.lists[m.focused].FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.state == stateDefault {
				m.focused = (m.focused + 1) % 2
				if m.focused == TaskListsPane {
					m.SetStatus("Task Lists")
				} else {
					m.SetStatus("Tasks")
				}
			}
		case "n":
			if m.state == stateDefault && m.focused == TaskListsPane {
				m.state = stateNewTaskList
				m.SetStatus("New Task List")
			}
		case "d":
			if m.state == stateDefault {
				if m.focused == TaskListsPane {
					m.state = stateDeleteTaskList
					m.SetStatus("Delete Task List? (y/n)")
				} else {
					m.state = stateDeleteTask
					m.SetStatus("Delete Task? (y/n)")
				}
			}
		case "y":
			if m.state == stateDeleteTaskList {
				m.state = stateDefault
				m.SetStatus("Deleting task list...")
				selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
				return m, func() tea.Msg {
					err := m.client.DeleteTaskList(gtasks.DeleteTaskListOptions{TaskListID: selectedTaskList.Id})
					if err != nil {
						return errorMsg{err}
					}
					return m.Init()
				}
			} else if m.state == stateDeleteTask {
				m.state = stateDefault
				m.SetStatus("Deleting task...")
				selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
				selectedTask := m.lists[TasksPane].SelectedItem().(taskItem)
				return m, func() tea.Msg {
					err := m.client.DeleteTask(gtasks.DeleteTaskOptions{TaskListID: selectedTaskList.Id, TaskID: selectedTask.Id})
					if err != nil {
						return errorMsg{err}
					}
					return taskDeletedMsg{}
				}
			}
		case "enter":
			if m.focused == TaskListsPane {
				selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
				return m, func() tea.Msg {
					tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{TaskListID: selectedTaskList.Id})
					if err != nil {
						return errorMsg{err}
					}
					return tasksLoadedMsg{tasks: tasks}
				}
			}
		case "esc":
			if m.state == stateNewTaskList || m.state == stateDeleteTaskList {
				m.state = stateDefault
				m.SetStatus("Ready")
			}
		}
	}

	var cmd tea.Cmd
	if m.state == stateNewTaskList {
		m.newTaskListInput, cmd = m.newTaskListInput.Update(msg)
	} else {
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	}
	return m, cmd
}