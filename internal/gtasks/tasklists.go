package gtasks

import "google.golang.org/api/tasks/v1"

// GetTaskListOptions holds the parameters for retrieving a task list.
type GetTaskListOptions struct {
	TaskListID string
}

// CreateTaskListOptions holds the parameters for creating a new task list.
type CreateTaskListOptions struct {
	Title string
}

// UpdateTaskListOptions holds the parameters for updating a task list.
type UpdateTaskListOptions struct {
	TaskListID string
	Title      string
}

// DeleteTaskListOptions holds the parameters for deleting a task list.
type DeleteTaskListOptions struct {
	TaskListID string
}

func (c *onlineClient) ListTaskLists() (*tasks.TaskLists, error) {
	return c.service.Tasklists.List().Do()
}

func (c *onlineClient) GetTaskList(opts GetTaskListOptions) (*tasks.TaskList, error) {
	return c.service.Tasklists.Get(opts.TaskListID).Do()
}

func (c *onlineClient) CreateTaskList(opts CreateTaskListOptions) (*tasks.TaskList, error) {
	list := &tasks.TaskList{
		Title: opts.Title,
	}
	return c.service.Tasklists.Insert(list).Do()
}

func (c *onlineClient) UpdateTaskList(opts UpdateTaskListOptions) (*tasks.TaskList, error) {
	list := &tasks.TaskList{
		Title: opts.Title,
	}
	return c.service.Tasklists.Update(opts.TaskListID, list).Do()
}

func (c *onlineClient) DeleteTaskList(opts DeleteTaskListOptions) error {
	return c.service.Tasklists.Delete(opts.TaskListID).Do()
}
