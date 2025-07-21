package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/gtasks"
)

var tasksCmd = &cobra.Command{
	Use:     "tasks",
	Short:   "Manage your tasks",
	Aliases: []string{"t"},
}

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks in a task list",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		showCompleted, _ := cmd.Flags().GetBool("show-completed")
		showHidden, _ := cmd.Flags().GetBool("show-hidden")

		opts := gtasks.ListTasksOptions{
			TaskListID:    tasklist,
			ShowCompleted: showCompleted,
			ShowHidden:    showHidden,
		}

		if err := client.ListTasks(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing tasks: %v\n", err)
			os.Exit(1)
		}
	},
}

var getTaskCmd = &cobra.Command{
	Use:   "get [ID]",
	Short: "Get details for a specific task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.GetTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		if err := client.GetTask(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error getting task: %v\n", err)
			os.Exit(1)
		}
	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		title, _ := cmd.Flags().GetString("title")
		notes, _ := cmd.Flags().GetString("notes")
		due, _ := cmd.Flags().GetString("due")

		opts := gtasks.CreateTaskOptions{
			TaskListID: tasklist,
			Title:      title,
			Notes:      notes,
			Due:        due,
		}

		if err := client.CreateTask(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating task: %v\n", err)
			os.Exit(1)
		}
	},
}

var updateTaskCmd = &cobra.Command{
	Use:   "update [ID]",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		title, _ := cmd.Flags().GetString("title")
		notes, _ := cmd.Flags().GetString("notes")
		due, _ := cmd.Flags().GetString("due")

		opts := gtasks.UpdateTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
			Title:      title,
			Notes:      notes,
			Due:        due,
		}

		if err := client.UpdateTask(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating task: %v\n", err)
			os.Exit(1)
		}
	},
}

var completeTaskCmd = &cobra.Command{
	Use:   "complete [ID]",
	Short: "Mark a task as complete",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.CompleteTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		if err := client.CompleteTask(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error completing task: %v\n", err)
			os.Exit(1)
		}
	},
}

var deleteTaskCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.DeleteTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		if err := client.DeleteTask(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting task: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tasksCmd)
	tasksCmd.AddCommand(listTasksCmd)
	tasksCmd.AddCommand(getTaskCmd)
	tasksCmd.AddCommand(createTaskCmd)
	tasksCmd.AddCommand(updateTaskCmd)
	tasksCmd.AddCommand(completeTaskCmd)
	tasksCmd.AddCommand(deleteTaskCmd)

	listTasksCmd.Flags().String("tasklist", "@default", "The ID of the task list")
	listTasksCmd.Flags().Bool("show-completed", false, "Include completed tasks in the output")
	listTasksCmd.Flags().Bool("show-hidden", false, "Include hidden tasks in the output")

	getTaskCmd.Flags().String("tasklist", "@default", "The ID of the task list")

	createTaskCmd.Flags().String("tasklist", "@default", "The ID of the task list")
	createTaskCmd.Flags().String("title", "", "The title of the new task")
	createTaskCmd.MarkFlagRequired("title")
	createTaskCmd.Flags().String("notes", "", "The notes for the new task")
	createTaskCmd.Flags().String("due", "", "The due date for the new task (RFC3339 format)")

	updateTaskCmd.Flags().String("tasklist", "@default", "The ID of the task list")
	updateTaskCmd.Flags().String("title", "", "The new title for the task")
	updateTaskCmd.Flags().String("notes", "", "The new notes for the task")
	updateTaskCmd.Flags().String("due", "", "The new due date for the task (RFC3339 format)")

	completeTaskCmd.Flags().String("tasklist", "@default", "The ID of the task list")

	deleteTaskCmd.Flags().String("tasklist", "@default", "The ID of the task list")
}