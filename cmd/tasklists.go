package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/gtasks"
	"github.com/yanicksenn/gtasks/internal/ui"
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

		lists, err := h.Client.ListTaskLists()
		if err != nil {
			return fmt.Errorf("error listing task lists: %w", err)
		}

		if !h.Quiet {
			if len(lists.Items) == 0 {
				cmd.Println("No task lists found.")
				return nil
			}

			cmd.Println("Task Lists:")
			for _, item := range lists.Items {
				cmd.Printf("- %s (%s)\n", item.Title, item.Id)
			}
		}
		return nil
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

		opts := gtasks.GetTaskListOptions{
			TaskListID: args[0],
		}

		list, err := h.Client.GetTaskList(opts)
		if err != nil {
			return fmt.Errorf("error getting task list: %w", err)
		}

		if !h.Quiet {
			cmd.Printf("ID:    %s\n", list.Id)
			cmd.Printf("Title: %s\n", list.Title)
			cmd.Printf("Self:  %s\n", list.SelfLink)
		}
		return nil
	},
}

var createTasklistCmd = &cobra.Command{
	Use:	 "create",
	Short:	 "Create a new task list",
	RunE: func(cmd *cobra.Command, args []string) error {
		return CreateTaskList(cmd, args, cmd.OutOrStdout())
	},
}

func CreateTaskList(cmd *cobra.Command, args []string, out io.Writer) error {
	h, err := NewCommandHelper(cmd)
	if err != nil {
		return err
	}

	title, _ := cmd.Flags().GetString("title")
	opts := gtasks.CreateTaskListOptions{
		Title: title,
	}

	createdList, err := h.Client.CreateTaskList(opts)
	if err != nil {
		return fmt.Errorf("error creating task list: %w", err)
	}

	if !h.Quiet {
		fmt.Fprintf(out, "Successfully created task list: %s (%s)\n", createdList.Title, createdList.Id)
	}
	return nil
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

		title, _ := cmd.Flags().GetString("title")
	opts := gtasks.UpdateTaskListOptions{
		TaskListID: args[0],
		Title:	    title,
	}

	updatedList, err := h.Client.UpdateTaskList(opts)
	if err != nil {
		return fmt.Errorf("error updating task list: %w", err)
	}

	if !h.Quiet {
		fmt.Printf("Successfully updated task list: %s (%s)\n", updatedList.Title, updatedList.Id)
	}
	return nil
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

		opts := gtasks.DeleteTaskListOptions{
			TaskListID: args[0],
		}

		err = h.Client.DeleteTaskList(opts)
		if err != nil {
			return fmt.Errorf("error deleting task list: %w", err)
		}

		if !h.Quiet {
			fmt.Printf("Successfully deleted task list: %s\n", args[0])
		}
		return nil
	},
}

var printTaskListCmd = &cobra.Command{
	Use:	 "print [ID]",
	Short:	 "Print a property of a task list",
	Long: `Print a property of a task list.
Available properties: id, title, selfLink`,
	Args:	 cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		opts := gtasks.GetTaskListOptions{
			TaskListID: args[0],
		}

		list, err := h.Client.GetTaskList(opts)
		if err != nil {
			return fmt.Errorf("error getting task list: %w", err)
		}

		property, _ := cmd.Flags().GetString("property")
		if err := ui.PrintTaskListProperty(list, property, h.Quiet); err != nil {
			return fmt.Errorf("error printing property: %w", err)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(tasklistsCmd)
	tasklistsCmd.AddCommand(listTasklistsCmd)
	tasklistsCmd.AddCommand(getTasklistCmd)
	tasklistsCmd.AddCommand(createTasklistCmd)
	tasklistsCmd.AddCommand(updateTasklistCmd)
	tasklistsCmd.AddCommand(deleteTasklistCmd)
	tasklistsCmd.AddCommand(printTaskListCmd)

	createTasklistCmd.Flags().String("title", "", "The title of the new task list")
	createTasklistCmd.MarkFlagRequired("title")

	updateTasklistCmd.Flags().String("title", "", "The new title for the task list")
	updateTasklistCmd.MarkFlagRequired("title")

	printTaskListCmd.Flags().String("property", "", "The property to print")
	printTaskListCmd.MarkFlagRequired("property")
}
