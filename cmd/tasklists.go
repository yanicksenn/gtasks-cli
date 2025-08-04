package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/gtasks"
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
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		lists, err := client.ListTaskLists()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing task lists: %v\n", err)
			os.Exit(1)
		}

		quiet, _ := cmd.Flags().GetBool("quiet")
		if !quiet {
			if len(lists.Items) == 0 {
				fmt.Println("No task lists found.")
				return
			}

			fmt.Println("Task Lists:")
			for _, item := range lists.Items {
				fmt.Printf("- %s (%s)\n", item.Title, item.Id)
			}
		}
	},
}

var getTasklistCmd = &cobra.Command{
	Use:   "get [ID]",
	Short: "Get details for a specific task list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		opts := gtasks.GetTaskListOptions{
			TaskListID: args[0],
		}

		list, err := client.GetTaskList(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting task list: %v\n", err)
			os.Exit(1)
		}

		quiet, _ := cmd.Flags().GetBool("quiet")
		if !quiet {
			fmt.Printf("ID:    %s\n", list.Id)
			fmt.Printf("Title: %s\n", list.Title)
			fmt.Printf("Self:  %s\n", list.SelfLink)
		}
	},
}

var createTasklistCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task list",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		title, _ := cmd.Flags().GetString("title")
		opts := gtasks.CreateTaskListOptions{
			Title: title,
		}

		createdList, err := client.CreateTaskList(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating task list: %v\n", err)
			os.Exit(1)
		}

		quiet, _ := cmd.Flags().GetBool("quiet")
		if !quiet {
			fmt.Printf("Successfully created task list: %s (%s)\n", createdList.Title, createdList.Id)
		}
	},
}

var updateTasklistCmd = &cobra.Command{
	Use:   "update [ID]",
	Short: "Update a task list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		title, _ := cmd.Flags().GetString("title")
		opts := gtasks.UpdateTaskListOptions{
			TaskListID: args[0],
			Title:      title,
		}

		updatedList, err := client.UpdateTaskList(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating task list: %v\n", err)
			os.Exit(1)
		}

		quiet, _ := cmd.Flags().GetBool("quiet")
		if !quiet {
			fmt.Printf("Successfully updated task list: %s (%s)\n", updatedList.Title, updatedList.Id)
		}
	},
}

var deleteTasklistCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Delete a task list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		opts := gtasks.DeleteTaskListOptions{
			TaskListID: args[0],
		}

		err = client.DeleteTaskList(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting task list: %v\n", err)
			os.Exit(1)
		}

		quiet, _ := cmd.Flags().GetBool("quiet")
		if !quiet {
			fmt.Printf("Successfully deleted task list: %s\n", args[0])
		}
	},
}

var printTaskListCmd = &cobra.Command{
	Use:   "print [ID]",
	Short: "Print a property of a task list",
	Long: `Print a property of a task list.
Available properties: id, title, selfLink`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		opts := gtasks.GetTaskListOptions{
			TaskListID: args[0],
		}

		list, err := client.GetTaskList(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting task list: %v\n", err)
			os.Exit(1)
		}

		property, _ := cmd.Flags().GetString("property")
        quiet, _ := cmd.Flags().GetBool("quiet")
        if err := gtasks.PrintTaskListProperty(list, property, quiet); err != nil {
            fmt.Fprintf(os.Stderr, "Error printing property: %v\n", err)
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
	tasklistsCmd.AddCommand(printTaskListCmd)

	createTasklistCmd.Flags().String("title", "", "The title of the new task list")
	createTasklistCmd.MarkFlagRequired("title")

	updateTasklistCmd.Flags().String("title", "", "The new title for the task list")
	updateTasklistCmd.MarkFlagRequired("title")

	printTaskListCmd.Flags().String("property", "", "The property to print")
	printTaskListCmd.MarkFlagRequired("property")
}
