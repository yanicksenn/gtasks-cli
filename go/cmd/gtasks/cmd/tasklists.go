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

func init() {
	rootCmd.AddCommand(tasklistsCmd)
	tasklistsCmd.AddCommand(listTasklistsCmd)
	tasklistsCmd.AddCommand(getTasklistCmd)
}
