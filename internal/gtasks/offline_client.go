package gtasks

import (
	"github.com/yanicksenn/gtasks/internal/store"
	"google.golang.org/api/tasks/v1"
)

// offlineClient is a client that interacts with the offline store.
type offlineClient struct {
	store *store.InMemoryStore
}

// newOfflineClient creates a new client that works with the local offline store.
func newOfflineClient() (*offlineClient, error) {
	path, err := store.GetOfflineStorePath()
	if err != nil {
		return nil, err
	}
	s, err := store.NewInMemoryStore(path)
	if err != nil {
		return nil, err
	}
	return &offlineClient{store: s}, nil
}

func (c *offlineClient) ListTaskLists() (*tasks.TaskLists, error) {
	lists, err := c.store.ListTaskLists()
	if err != nil {
		return nil, err
	}
	return &tasks.TaskLists{Items: lists}, nil
}

func (c *offlineClient) CreateTaskList(opts CreateTaskListOptions) (*tasks.TaskList, error) {
	list := &tasks.TaskList{
		Title: opts.Title,
	}
	return c.store.CreateTaskList(list)
}

func (c *offlineClient) GetTaskList(opts GetTaskListOptions) (*tasks.TaskList, error) {
	return c.store.GetTaskList(opts.TaskListID)
}

func (c *offlineClient) UpdateTaskList(opts UpdateTaskListOptions) (*tasks.TaskList, error) {
	list := &tasks.TaskList{
		Title: opts.Title,
	}
	return c.store.UpdateTaskList(opts.TaskListID, list)
}

func (c *offlineClient) DeleteTaskList(opts DeleteTaskListOptions) error {
	return c.store.DeleteTaskList(opts.TaskListID)
}

func (c *offlineClient) ListTasks(opts ListTasksOptions) (*tasks.Tasks, error) {
	taskItems, err := c.store.ListTasks(opts.TaskListID)
	if err != nil {
		return nil, err
	}
	return &tasks.Tasks{Items: taskItems}, nil
}

func (c *offlineClient) CreateTask(opts CreateTaskOptions) (*tasks.Task, error) {
	task := &tasks.Task{
		Title: opts.Title,
		Notes: opts.Notes,
		Due:   opts.Due,
	}
	return c.store.CreateTask(opts.TaskListID, task)
}

func (c *offlineClient) GetTask(opts GetTaskOptions) (*tasks.Task, error) {
	return c.store.GetTask(opts.TaskListID, opts.TaskID)
}

func (c *offlineClient) UpdateTask(opts UpdateTaskOptions) (*tasks.Task, error) {
	task := &tasks.Task{
		Title: opts.Title,
		Notes: opts.Notes,
		Due:   opts.Due,
	}
	return c.store.UpdateTask(opts.TaskListID, opts.TaskID, task)
}

func (c *offlineClient) CompleteTask(opts CompleteTaskOptions) (*tasks.Task, error) {
	task, err := c.store.GetTask(opts.TaskListID, opts.TaskID)
	if err != nil {
		return nil, err
	}
	task.Status = "completed"
	return c.store.UpdateTask(opts.TaskListID, opts.TaskID, task)
}

func (c *offlineClient) UncompleteTask(opts UncompleteTaskOptions) (*tasks.Task, error) {
	task, err := c.store.GetTask(opts.TaskListID, opts.TaskID)
	if err != nil {
		return nil, err
	}
	task.Status = "needsAction"
	return c.store.UpdateTask(opts.TaskListID, opts.TaskID, task)
}

func (c *offlineClient) DeleteTask(opts DeleteTaskOptions) error {
	return c.store.DeleteTask(opts.TaskListID, opts.TaskID)
}
