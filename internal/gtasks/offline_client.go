package gtasks

import (
	"sort"
	"strings"

	"github.com/yanicksenn/gtasks/internal/store"
	"google.golang.org/api/tasks/v1"
)

// offlineClient is a client that interacts with the offline store.
// It implements the Client interface.
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

// ListTaskLists lists the task lists from the offline store.
func (c *offlineClient) ListTaskLists(opts ListTaskListsOptions) (*tasks.TaskLists, error) {
	lists, err := c.store.ListTaskLists()
	if err != nil {
		return nil, err
	}

	switch opts.SortBy {
	case "alphabetical":
		sort.Slice(lists, func(i, j int) bool {
			return strings.ToLower(lists[i].Title) < strings.ToLower(lists[j].Title)
		})
	case "last-modified":
		sort.Slice(lists, func(i, j int) bool {
			return lists[i].Updated > lists[j].Updated
		})
	case "uncompleted-tasks":
		// This is a placeholder. The actual implementation will require fetching tasks for each list.
	}

	return &tasks.TaskLists{Items: lists}, nil
}

// CreateTaskList creates a new task list in the offline store.
func (c *offlineClient) CreateTaskList(opts CreateTaskListOptions) (*tasks.TaskList, error) {
	list := &tasks.TaskList{
		Title: opts.Title,
	}
	return c.store.CreateTaskList(list)
}

// GetTaskList retrieves a task list from the offline store.
func (c *offlineClient) GetTaskList(opts GetTaskListOptions) (*tasks.TaskList, error) {
	return c.store.GetTaskList(opts.TaskListID)
}

// UpdateTaskList updates a task list in the offline store.
func (c *offlineClient) UpdateTaskList(opts UpdateTaskListOptions) (*tasks.TaskList, error) {
	list := &tasks.TaskList{
		Title: opts.Title,
	}
	return c.store.UpdateTaskList(opts.TaskListID, list)
}

// DeleteTaskList deletes a task list from the offline store.
func (c *offlineClient) DeleteTaskList(opts DeleteTaskListOptions) error {
	return c.store.DeleteTaskList(opts.TaskListID)
}

// ListTasks lists the tasks from the offline store.
func (c *offlineClient) ListTasks(opts ListTasksOptions) (*tasks.Tasks, error) {
	taskItems, err := c.store.ListTasks(opts.TaskListID)
	if err != nil {
		return nil, err
	}
	return &tasks.Tasks{Items: taskItems}, nil
}

// CreateTask creates a new task in the offline store.
func (c *offlineClient) CreateTask(opts CreateTaskOptions) (*tasks.Task, error) {
	task := &tasks.Task{
		Title: opts.Title,
		Notes: opts.Notes,
		Due:   opts.Due,
	}
	return c.store.CreateTask(opts.TaskListID, task)
}

// GetTask retrieves a task from the offline store.
func (c *offlineClient) GetTask(opts GetTaskOptions) (*tasks.Task, error) {
	return c.store.GetTask(opts.TaskListID, opts.TaskID)
}

// UpdateTask updates a task in the offline store.
func (c *offlineClient) UpdateTask(opts UpdateTaskOptions) (*tasks.Task, error) {
	task := &tasks.Task{
		Title: opts.Title,
		Notes: opts.Notes,
		Due:   opts.Due,
	}
	return c.store.UpdateTask(opts.TaskListID, opts.TaskID, task)
}

// CompleteTask marks a task as complete in the offline store.
func (c *offlineClient) CompleteTask(opts CompleteTaskOptions) (*tasks.Task, error) {
	task, err := c.store.GetTask(opts.TaskListID, opts.TaskID)
	if err != nil {
		return nil, err
	}
	task.Status = "completed"
	return c.store.UpdateTask(opts.TaskListID, opts.TaskID, task)
}

// UncompleteTask marks a task as not complete in the offline store.
func (c *offlineClient) UncompleteTask(opts UncompleteTaskOptions) (*tasks.Task, error) {
	task, err := c.store.GetTask(opts.TaskListID, opts.TaskID)
	if err != nil {
		return nil, err
	}
	task.Status = "needsAction"
	return c.store.UpdateTask(opts.TaskListID, opts.TaskID, task)
}

// DeleteTask deletes a task from the offline store.
func (c *offlineClient) DeleteTask(opts DeleteTaskOptions) error {
	return c.store.DeleteTask(opts.TaskListID, opts.TaskID)
}
