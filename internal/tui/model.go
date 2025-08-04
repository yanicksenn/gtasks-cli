package tui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanicksenn/gtasks/internal/gtasks"
	taskspb "google.golang.org/api/tasks/v1"
)

type tasksFetchTimeoutMsg struct{}

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

type taskCompletedMsg struct{}

type taskUncompletedMsg struct{}

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
	client           gtasks.Client
	focused          Pane
	lists            []list.Model
	status           string
	state            state
	newTaskListInput textinput.Model
	sortBy           []string
	timer            *time.Timer
	delegate         *itemDelegate
	keys             keyMap
	width            int
	height           int
	detailsViewHeight int
}

func New(offline bool) (*Model, error) {
	client, err := gtasks.NewClient(context.Background(), offline)
	if err != nil {
		return nil, err
	}

	newTaskListInput := textinput.New()
	newTaskListInput.Placeholder = "New Task List"
	newTaskListInput.Focus()

	delegate := &itemDelegate{focused: true}
	taskLists := list.New([]list.Item{}, delegate, 0, 0)
	taskLists.Title = "Task Lists"
	taskLists.SetShowHelp(false)
	tasks := list.New([]list.Item{}, taskItemDelegate{}, 0, 0)
	tasks.Title = "Tasks"
	tasks.SetShowHelp(false)

	m := &Model{
		client:           client,
		focused:          TaskListsPane,
		lists:            []list.Model{taskLists, tasks},
		state:            stateDefault,
		newTaskListInput: newTaskListInput,
		sortBy:           []string{"alphabetical", "last-modified", "uncompleted-tasks"},
		timer:            time.NewTimer(0),
		delegate:         delegate,
		keys:             keys,
		detailsViewHeight: 10,
	}
	m.timer.Stop()
	m.SetStatus("Ready")
	return m, nil
}

func (m *Model) sort() {
	if m.focused == TaskListsPane {
		m.sortBy = append(m.sortBy[1:], m.sortBy[0])
		m.SetStatus("Sorted by " + m.sortBy[0])
		m.Init()
	} else {
		m.sortBy = append(m.sortBy[1:], m.sortBy[0])
		m.SetStatus("Sorted by " + m.sortBy[0])
		m.Update(tasksLoadedMsg{tasks: &taskspb.Tasks{}})
		m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	}
}


func (m *Model) SetStatus(status string) {
	m.status = status
}

func (m *Model) Init() tea.Cmd {
	return func() tea.Msg {
		taskLists, err := m.client.ListTaskLists(gtasks.ListTaskListsOptions{SortBy: m.sortBy[0]})
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
		m.width = msg.Width
		m.height = msg.Height
		listHeight := m.height - m.detailsViewHeight - v - 1
		m.lists[TaskListsPane].SetSize(m.width/2-h, listHeight)
		m.lists[TasksPane].SetSize(m.width/2-h, listHeight)
		return m, nil

	case errorMsg:
		m.SetStatus(msg.Error())
		return m, nil

	case tea.KeyMsg:
		return m.handleKeypress(msg)

	default:
		return m.handleOther(msg)
	}
}

func (m *Model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.lists[m.focused].FilterState() == list.Filtering {
		goto update
	}

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Sort):
		return m.handleSort(msg)
	case key.Matches(msg, m.keys.New):
		return m.handleNew(msg)
	case key.Matches(msg, m.keys.Delete):
		return m.handleDelete(msg)
	case key.Matches(msg, m.keys.Confirm):
		return m.handleConfirm(msg)
	case key.Matches(msg, m.keys.Select):
		return m.handleSelect(msg)
	case key.Matches(msg, m.keys.Back):
		return m.handleBack(msg)
	case key.Matches(msg, m.keys.Cancel):
		return m.handleCancel(msg)
	case key.Matches(msg, m.keys.ToggleComplete):
		return m.handleToggleComplete(msg)
	case key.Matches(msg, m.keys.Left):
		return m.handleLeft(msg)
	case key.Matches(msg, m.keys.Right):
		return m.handleRight(msg)
	case key.Matches(msg, m.keys.Up, m.keys.Down):
		return m.handleNavigation(msg)
	}

update:
	var cmd tea.Cmd
	if m.state == stateNewTaskList {
		m.newTaskListInput, cmd = m.newTaskListInput.Update(msg)
	} else {
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	}
	return m, cmd
}

func (m *Model) handleOther(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case taskListsLoadedMsg:
		items := make([]list.Item, len(msg.taskLists.Items))
		for i, taskList := range msg.taskLists.Items {
			items[i] = taskListItem{taskList}
		}
		m.lists[TaskListsPane].SetItems(items)
		return m, nil

	case tasksLoadedMsg:
		items := make([]list.Item, len(msg.tasks.Items))
		for i, task := range msg.tasks.Items {
			items[i] = taskItem{task}
		}
		m.lists[TasksPane].SetItems(items)
		if len(items) == 0 {
			m.lists[TasksPane].Title = "Empty task list"
		} else {
			m.lists[TasksPane].Title = "Tasks"
		}
		return m, nil

	case tasksFetchTimeoutMsg:
		selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
		return m, func() tea.Msg {
			tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{TaskListID: selectedTaskList.Id, ShowCompleted: true, SortBy: m.sortBy[0]})
			if err != nil {
				return errorMsg{err}
			}
			return tasksLoadedMsg{tasks: tasks}
		}

	case taskListCreatedMsg:
		m.newTaskListInput.Reset()
		m.state = stateDefault
		m.SetStatus("Task list created")
		return m, m.Init()

	case taskDeletedMsg:
		m.SetStatus("Task deleted")
		selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
		return m, func() tea.Msg {
			tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{TaskListID: selectedTaskList.Id, ShowCompleted: true, SortBy: m.sortBy[0]})
			if err != nil {
				return errorMsg{err}
			}
			return tasksLoadedMsg{tasks: tasks}
		}

	case taskCompletedMsg:
		m.SetStatus("Task completed")
		selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
		return m, func() tea.Msg {
			tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{TaskListID: selectedTaskList.Id, ShowCompleted: true, SortBy: m.sortBy[0]})
			if err != nil {
				return errorMsg{err}
			}
			return tasksLoadedMsg{tasks: tasks}
		}

	case taskUncompletedMsg:
		m.SetStatus("Task un-completed")
		selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
		return m, func() tea.Msg {
			tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{TaskListID: selectedTaskList.Id, ShowCompleted: true, SortBy: m.sortBy[0]})
			if err != nil {
				return errorMsg{err}
			}
			return tasksLoadedMsg{tasks: tasks}
		}
	}
	return m, nil
}

func (m *Model) handleSort(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == stateDefault {
		m.sort()
	}
	return m, nil
}

func (m *Model) handleTab(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == stateDefault {
		m.focused = (m.focused + 1) % 2
		if m.focused == TaskListsPane {
			m.delegate.focused = true
			m.SetStatus("Task Lists")
		} else {
			m.delegate.focused = false
			m.SetStatus("Tasks")
		}
	}
	return m, nil
}

func (m *Model) handleNew(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == stateDefault && m.focused == TaskListsPane {
		m.state = stateNewTaskList
		m.SetStatus("New Task List")
	}
	return m, nil
}

func (m *Model) handleDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == stateDefault {
		if m.focused == TaskListsPane {
			m.state = stateDeleteTaskList
			m.SetStatus("Delete Task List? (y/n)")
		} else {
			m.state = stateDeleteTask
			m.SetStatus("Delete Task? (y/n)")
		}
	}
	return m, nil
}

func (m *Model) handleConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	return m, nil
}

func (m *Model) handleSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == stateNewTaskList {
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
	}

	if m.focused == TaskListsPane {
		m.focused = TasksPane
		m.SetStatus("Tasks")
		selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
		return m, func() tea.Msg {
			tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{TaskListID: selectedTaskList.Id, ShowCompleted: true, SortBy: m.sortBy[0]})
			if err != nil {
				return errorMsg{err}
			}
			return tasksLoadedMsg{tasks: tasks}
		}
	} else if m.focused == TasksPane {
		// No-op, selection is handled by the list component
	}
	return m, nil
}

func (m *Model) handleBack(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) handleCancel(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == stateNewTaskList {
		m.state = stateDefault
		m.SetStatus("Ready")
		return m, nil
	}

	if m.state == stateDeleteTaskList || m.state == stateDeleteTask {
		m.state = stateDefault
		m.SetStatus("Ready")
		return m, nil
	}
	return m, nil
}

func (m *Model) handleToggleComplete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == stateDefault && m.focused == TasksPane {
		selectedTask := m.lists[TasksPane].SelectedItem().(taskItem)
		if selectedTask.Status == "completed" {
			m.SetStatus("Un-completing task...")
			selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
			return m, func() tea.Msg {
				_, err := m.client.UncompleteTask(gtasks.UncompleteTaskOptions{TaskListID: selectedTaskList.Id, TaskID: selectedTask.Id})
				if err != nil {
					return errorMsg{err}
				}
				return taskUncompletedMsg{}
			}
		} else {
			m.SetStatus("Completing task...")
			selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
			return m, func() tea.Msg {
				_, err := m.client.CompleteTask(gtasks.CompleteTaskOptions{TaskListID: selectedTaskList.Id, TaskID: selectedTask.Id})
				if err != nil {
					return errorMsg{err}
				}
				return taskCompletedMsg{}
			}
		}
	}
	return m, nil
}

func (m *Model) handleLeft(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.focused == TasksPane {
		m.focused = TaskListsPane
		m.SetStatus("Task Lists")
		return m, nil
	}
	return m, nil
}

func (m *Model) handleRight(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.focused == TaskListsPane {
		m.focused = TasksPane
		m.SetStatus("Tasks")
		selectedTaskList := m.lists[TaskListsPane].SelectedItem().(taskListItem)
		return m, func() tea.Msg {
			tasks, err := m.client.ListTasks(gtasks.ListTasksOptions{TaskListID: selectedTaskList.Id, ShowCompleted: true, SortBy: m.sortBy[0]})
			if err != nil {
				return errorMsg{err}
			}
			return tasksLoadedMsg{tasks: tasks}
		}
	}
	return m, nil
}

func (m *Model) handleNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var timeoutCmd tea.Cmd
	if m.focused == TaskListsPane {
		m.timer.Reset(300 * time.Millisecond)
		m.lists[TasksPane].Title = "Loading..."
		m.lists[TasksPane].SetItems([]list.Item{})
		timeoutCmd = func() tea.Msg {
			<-m.timer.C
			return tasksFetchTimeoutMsg{}
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, tea.Batch(cmd, timeoutCmd)
}