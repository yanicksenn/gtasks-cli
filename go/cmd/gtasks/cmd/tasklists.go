package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/gtasks"
)

var tasklistsCmd = &cobra.Command{
	Use:     "tasklists",
	Short:   "Manage your task lists",
	Aliases: []string{"tl"},
}

var listTasklistsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your task lists",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		if err := client.ListTaskLists(); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing task lists: %v\n", err)
			os.Exit(1)
		}
	},
}

var getTasklistCmd = &cobra.Command{
	Use:   "get [ID]",
	Short: "Get details for a specific task list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		opts := gtasks.GetTaskListOptions{
			TaskListID: args[0],
		}

		if err := client.GetTaskList(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error getting task list: %v\n", err)
			os.Exit(1)
		}
	},
}

var createTasklistCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task list",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		title, _ := cmd.Flags().GetString("title")
		opts := gtasks.CreateTaskListOptions{
			Title: title,
		}

		if err := client.CreateTaskList(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating task list: %v\n", err)
			os.Exit(1)
		}
	},
}

var updateTasklistCmd = &cobra.Command{
	Use:   "update [ID]",
	Short: "Update a task list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		title, _ := cmd.Flags().GetString("title")
		opts := gtasks.UpdateTaskListOptions{
			TaskListID: args[0],
			Title:      title,
		}

		if err := client.UpdateTaskList(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating task list: %v\n", err)
			os.Exit(1)
		}
	},
}

var deleteTasklistCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Delete a task list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		opts := gtasks.DeleteTaskListOptions{
			TaskListID: args[0],
		}

		if err := client.DeleteTaskList(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting task list: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tasklistsCmd)
	tasklistsCmd.AddCommand(listTasklistsCmd)
	tasklistsCmd.AddCommand(getTasklistCmd)
	tasklistsCmd.AddCommand(createTasklistCmd)
	tasklistsCmd.AddCommand(updateTasklistCmd)
	tasklistsCmd.AddCommand(deleteTasklistCmd)

	createTasklistCmd.Flags().String("title", "", "The title of the new task list")
	createTasklistCmd.MarkFlagRequired("title")

	updateTasklistCmd.Flags().String("title", "", "The new title for the task list")
	updateTasklistCmd.MarkFlagRequired("title")
}