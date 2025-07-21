package gtasks

import (
	"fmt"

	"google.golang.org/api/tasks/v1"
)

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

// CreateTaskList creates a new task list.
func (c *Client) CreateTaskList(opts CreateTaskListOptions) error {
	list := &tasks.TaskList{
		Title: opts.Title,
	}

	createdList, err := c.service.Tasklists().Insert(list).Do()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully created task list: %s (%s)\n", createdList.Title, createdList.Id)
	return nil
}

// UpdateTaskList updates a task list.
func (c *Client) UpdateTaskList(opts UpdateTaskListOptions) error {
	list := &tasks.TaskList{
		Title: opts.Title,
	}

	updatedList, err := c.service.Tasklists().Update(opts.TaskListID, list).Do()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully updated task list: %s (%s)\n", updatedList.Title, updatedList.Id)
	return nil
}

// DeleteTaskList deletes a task list.
func (c *Client) DeleteTaskList(opts DeleteTaskListOptions) error {
	err := c.service.Tasklists().Delete(opts.TaskListID).Do()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully deleted task list: %s\n", opts.TaskListID)
	return nil
}