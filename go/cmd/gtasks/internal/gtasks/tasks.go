package gtasks

import (
	"fmt"

	"google.golang.org/api/tasks/v1"
)

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

// ListTasks lists all tasks in a task list.
func (c *Client) ListTasks(opts ListTasksOptions) error {
	tasks, err := c.service.Tasks().List(opts.TaskListID).ShowCompleted(opts.ShowCompleted).ShowHidden(opts.ShowHidden).Do()
	if err != nil {
		return err
	}

	if len(tasks.Items) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	fmt.Println("Tasks:")
	for _, item := range tasks.Items {
		status := " "
		if item.Status == "completed" {
			status = "x"
		}
		fmt.Printf("[%s] %s (%s)\n", status, item.Title, item.Id)
	}

	return nil
}

// GetTask retrieves a single task.
func (c *Client) GetTask(opts GetTaskOptions) error {
	task, err := c.service.Tasks().Get(opts.TaskListID, opts.TaskID).Do()
	if err != nil {
		return err
	}

	fmt.Printf("ID:      %s\n", task.Id)
	fmt.Printf("Title:   %s\n", task.Title)
	fmt.Printf("Status:  %s\n", task.Status)
	fmt.Printf("Notes:   %s\n", task.Notes)
	fmt.Printf("Due:     %s\n", task.Due)
	fmt.Printf("Self:    %s\n", task.SelfLink)

	return nil
}

// CreateTask creates a new task.
func (c *Client) CreateTask(opts CreateTaskOptions) error {
	task := &tasks.Task{
		Title: opts.Title,
		Notes: opts.Notes,
		Due:   opts.Due,
	}

	createdTask, err := c.service.Tasks().Insert(opts.TaskListID, task).Do()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully created task: %s (%s)\n", createdTask.Title, createdTask.Id)
	return nil
}

// UpdateTask updates a task.
func (c *Client) UpdateTask(opts UpdateTaskOptions) error {
	task := &tasks.Task{
		Title: opts.Title,
		Notes: opts.Notes,
		Due:   opts.Due,
	}

	updatedTask, err := c.service.Tasks().Update(opts.TaskListID, opts.TaskID, task).Do()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully updated task: %s (%s)\n", updatedTask.Title, updatedTask.Id)
	return nil
}

// CompleteTask marks a task as complete.
func (c *Client) CompleteTask(opts CompleteTaskOptions) error {
	task, err := c.service.Tasks().Get(opts.TaskListID, opts.TaskID).Do()
	if err != nil {
		return err
	}

	task.Status = "completed"

	completedTask, err := c.service.Tasks().Update(opts.TaskListID, opts.TaskID, task).Do()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully completed task: %s (%s)\n", completedTask.Title, completedTask.Id)
	return nil
}

// DeleteTask deletes a task.
func (c *Client) DeleteTask(opts DeleteTaskOptions) error {
	err := c.service.Tasks().Delete(opts.TaskListID, opts.TaskID).Do()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully deleted task: %s\n", opts.TaskID)
	return nil
}