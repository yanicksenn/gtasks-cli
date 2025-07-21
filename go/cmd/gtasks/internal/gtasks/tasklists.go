package gtasks

import (
	"fmt"
)

// GetTaskListOptions holds the parameters for retrieving a task list.
type GetTaskListOptions struct {
	TaskListID string
}

// ListTaskLists lists all task lists.
func (c *Client) ListTaskLists() error {
	lists, err := c.service.Tasklists().List().Do()
	if err != nil {
		return err
	}

	if len(lists.Items) == 0 {
		fmt.Println("No task lists found.")
		return nil
	}

	fmt.Println("Task Lists:")
	for _, item := range lists.Items {
		fmt.Printf("- %s (%s)\n", item.Title, item.Id)
	}

	return nil
}

// GetTaskList retrieves a single task list.
func (c *Client) GetTaskList(opts GetTaskListOptions) error {
	list, err := c.service.Tasklists().Get(opts.TaskListID).Do()
	if err != nil {
		return err
	}

	fmt.Printf("ID:    %s\n", list.Id)
	fmt.Printf("Title: %s\n", list.Title)
	fmt.Printf("Self:  %s\n", list.SelfLink)

	return nil
}
