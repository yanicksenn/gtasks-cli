package gtasks

import (
	"fmt"
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
