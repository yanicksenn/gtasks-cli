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
		client, err := gtasks.NewClient(context.Background())
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
		client, err := gtasks.NewClient(context.Background())
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

func init() {
	rootCmd.AddCommand(tasksCmd)
	tasksCmd.AddCommand(listTasksCmd)
	tasksCmd.AddCommand(getTaskCmd)

	listTasksCmd.Flags().String("tasklist", "", "The ID of the task list")
	listTasksCmd.MarkFlagRequired("tasklist")
	listTasksCmd.Flags().Bool("show-completed", false, "Include completed tasks in the output")
	listTasksCmd.Flags().Bool("show-hidden", false, "Include hidden tasks in the output")

	getTaskCmd.Flags().String("tasklist", "", "The ID of the task list")
	getTaskCmd.MarkFlagRequired("tasklist")
}
