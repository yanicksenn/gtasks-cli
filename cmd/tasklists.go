package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/gtasks"
)

var tasklistsCmd = &cobra.Command{
	Use:	 "tasklists",
	Short:	 "Manage your task lists",
	Aliases: []string{"tl"},
}

var listTasklistsCmd = &cobra.Command{
	Use:	 "list",
	Short:	 "List all your task lists",
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		// Get the sort-by flag value
		sortBy, _ := cmd.Flags().GetString("sort-by")
		opts := gtasks.ListTaskListsOptions{
			SortBy: sortBy,
		}

		// List the task lists
		lists, err := h.Client.ListTaskLists(opts)
		if err != nil {
			return fmt.Errorf("error listing task lists: %w", err)
		}

		// Print the task lists
		return h.Printer.PrintTaskLists(lists)
	},
}

var getTasklistCmd = &cobra.Command{
	Use:	 "get [ID]",
	Short:	 "Get details for a specific task list",
	Args:	 cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		// Get the task list ID from the arguments
		opts := gtasks.GetTaskListOptions{
			TaskListID: args[0],
		}

		// Get the task list
		list, err := h.Client.GetTaskList(opts)
		if err != nil {
			return fmt.Errorf("error getting task list: %w", err)
		}

		// Print the task list
		return h.Printer.PrintTaskList(list)
	},
}

var createTasklistCmd = &cobra.Command{
	Use:	 "create",
	Short:	 "Create a new task list",
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		// Get the title from the flags
		title, _ := cmd.Flags().GetString("title")
		opts := gtasks.CreateTaskListOptions{
			Title: title,
		}

		// Create the task list
		createdList, err := h.Client.CreateTaskList(opts)
		if err != nil {
			return fmt.Errorf("error creating task list: %w", err)
		}

		// Print the created task list
		return h.Printer.PrintTaskList(createdList)
	},
}

var updateTasklistCmd = &cobra.Command{
	Use:	 "update [ID]",
	Short:	 "Update a task list",
	Args:	 cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		// Get the title from the flags
		title, _ := cmd.Flags().GetString("title")
		opts := gtasks.UpdateTaskListOptions{
			TaskListID: args[0],
			Title:	    title,
		}

		// Update the task list
		updatedList, err := h.Client.UpdateTaskList(opts)
		if err != nil {
			return fmt.Errorf("error updating task list: %w", err)
		}

		// Print a success message
		return h.Printer.PrintSuccess(fmt.Sprintf("Successfully updated task list: %s (%s)", updatedList.Title, updatedList.Id))
	},
}

var deleteTasklistCmd = &cobra.Command{
	Use:	 "delete [ID]",
	Short:	 "Delete a task list",
	Args:	 cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		// Get the task list ID from the arguments
		opts := gtasks.DeleteTaskListOptions{
			TaskListID: args[0],
		}

		// Delete the task list
		err = h.Client.DeleteTaskList(opts)
		if err != nil {
			return fmt.Errorf("error deleting task list: %w", err)
		}

		// Print a success message
		return h.Printer.PrintDelete("task list", args[0])
	},
}

func init() {
	RootCmd.AddCommand(tasklistsCmd)
	tasklistsCmd.AddCommand(listTasklistsCmd)
	tasklistsCmd.AddCommand(getTasklistCmd)
	tasklistsCmd.AddCommand(createTasklistCmd)
	tasklistsCmd.AddCommand(updateTasklistCmd)
	tasklistsCmd.AddCommand(deleteTasklistCmd)

	listTasklistsCmd.Flags().String("sort-by", "alphabetical", "Sort task lists by (alphabetical, last-modified, uncompleted-tasks)")

	createTasklistCmd.Flags().String("title", "", "The title of the new task list")
	createTasklistCmd.MarkFlagRequired("title")

	updateTasklistCmd.Flags().String("title", "", "The new title for the task list")
	updateTasklistCmd.MarkFlagRequired("title")
}
