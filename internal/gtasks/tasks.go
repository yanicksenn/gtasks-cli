package gtasks

import (
	"sort"
	"time"

	"google.golang.org/api/tasks/v1"
)

// ListTasksOptions holds the parameters for listing tasks.
type ListTasksOptions struct {
	TaskListID    string
	ShowCompleted bool
	ShowHidden    bool
	SortBy        string
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

// UncompleteTaskOptions holds the parameters for uncompleting a task.
type UncompleteTaskOptions struct {
	TaskListID string
	TaskID     string
}

// DeleteTaskOptions holds the parameters for deleting a task.
type DeleteTaskOptions struct {
	TaskListID string
	TaskID     string
}

func (c *onlineClient) ListTasks(opts ListTasksOptions) (*tasks.Tasks, error) {
	tasks, err := c.service.Tasks.List(opts.TaskListID).ShowCompleted(opts.ShowCompleted).ShowHidden(opts.ShowHidden).Do()
	if err != nil {
		return nil, err
	}

	switch opts.SortBy {
	case "alphabetical":
		sort.Slice(tasks.Items, func(i, j int) bool {
			return tasks.Items[i].Title < tasks.Items[j].Title
		})
	case "last-modified":
		sort.Slice(tasks.Items, func(i, j int) bool {
			return tasks.Items[i].Updated > tasks.Items[j].Updated
		})
	case "due-date":
		sort.Slice(tasks.Items, func(i, j int) bool {
			if tasks.Items[i].Due == "" {
				return false
			}
			if tasks.Items[j].Due == "" {
				return true
			}
			dueI, _ := time.Parse(time.RFC3339, tasks.Items[i].Due)
			dueJ, _ := time.Parse(time.RFC3339, tasks.Items[j].Due)
			if dueI.Equal(dueJ) {
				return tasks.Items[i].Title < tasks.Items[j].Title
			}
			return dueI.Before(dueJ)
		})
	}

	return tasks, nil
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

func (c *onlineClient) UncompleteTask(opts UncompleteTaskOptions) (*tasks.Task, error) {
	task, err := c.service.Tasks.Get(opts.TaskListID, opts.TaskID).Do()
	if err != nil {
		return nil, err
	}
	task.Status = "needsAction"
	return c.service.Tasks.Update(opts.TaskListID, opts.TaskID, task).Do()
}

func (c *onlineClient) DeleteTask(opts DeleteTaskOptions) error {
	return c.service.Tasks.Delete(opts.TaskListID, opts.TaskID).Do()
}
