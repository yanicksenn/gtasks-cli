package gtasks

import "google.golang.org/api/tasks/v1"

// ListTasksOptions holds the parameters for listing tasks.
type ListTasksOptions struct {
	TaskListID    string
	ShowCompleted bool
	ShowHidden    bool
}

// GetTaskOptions holds the parameters for retrieving a single task.
type GetTaskOptions struct {
	TaskListID string
	TaskID     string
}

// CreateTaskOptions holds the parameters for creating a new task.
type CreateTaskOptions struct {
	TaskListID string
	Title      string
	Notes      string
	Due        string
}

// UpdateTaskOptions holds the parameters for updating a task.
type UpdateTaskOptions struct {
	TaskListID string
	TaskID     string
	Title      string
	Notes      string
	Due        string
}

// CompleteTaskOptions holds the parameters for completing a task.
type CompleteTaskOptions struct {
	TaskListID string
	TaskID     string
}

// DeleteTaskOptions holds the parameters for deleting a task.
type DeleteTaskOptions struct {
	TaskListID string
	TaskID     string
}

func (c *onlineClient) ListTasks(opts ListTasksOptions) (*tasks.Tasks, error) {
	return c.service.Tasks.List(opts.TaskListID).ShowCompleted(opts.ShowCompleted).ShowHidden(opts.ShowHidden).Do()
}

func (c *onlineClient) GetTask(opts GetTaskOptions) (*tasks.Task, error) {
	return c.service.Tasks.Get(opts.TaskListID, opts.TaskID).Do()
}

func (c *onlineClient) CreateTask(opts CreateTaskOptions) (*tasks.Task, error) {
	task := &tasks.Task{
		Title: opts.Title,
		Notes: opts.Notes,
		Due:   opts.Due,
	}
	return c.service.Tasks.Insert(opts.TaskListID, task).Do()
}

func (c *onlineClient) UpdateTask(opts UpdateTaskOptions) (*tasks.Task, error) {
	task := &tasks.Task{
		Title: opts.Title,
		Notes: opts.Notes,
		Due:   opts.Due,
	}
	return c.service.Tasks.Update(opts.TaskListID, opts.TaskID, task).Do()
}

func (c *onlineClient) CompleteTask(opts CompleteTaskOptions) (*tasks.Task, error) {
	task, err := c.service.Tasks.Get(opts.TaskListID, opts.TaskID).Do()
	if err != nil {
		return nil, err
	}
	task.Status = "completed"
	return c.service.Tasks.Update(opts.TaskListID, opts.TaskID, task).Do()
}

func (c *onlineClient) DeleteTask(opts DeleteTaskOptions) error {
	return c.service.Tasks.Delete(opts.TaskListID, opts.TaskID).Do()
}
