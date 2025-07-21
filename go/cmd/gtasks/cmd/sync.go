package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/workspace/go/cmd/gtasks/internal/gtasks"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize offline data with Google Tasks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This will push all your local, offline changes to Google Tasks.")
		fmt.Println("WARNING: This operation will DELETE ALL existing online task lists and tasks and replace them with your local data.")
		fmt.Print("Are you sure you want to continue? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			fmt.Println("Sync cancelled.")
			return
		}

		// Create an offline client to read local data
		offlineCmd := &cobra.Command{}
		offlineCmd.Flags().Bool("offline", true, "")
		offlineClient, err := gtasks.NewClient(offlineCmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating offline client: %v\n", err)
			os.Exit(1)
		}

		// Create an online client to write to the API
		onlineCmd := &cobra.Command{}
		onlineCmd.Flags().Bool("offline", false, "")
		onlineClient, err := gtasks.NewClient(onlineCmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating online client: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Starting sync...")

		// 1. Get all offline data
		offlineTaskLists, err := offlineClient.ListTaskLists()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting offline task lists: %v\n", err)
			os.Exit(1)
		}

		// 2. Get all online data
		onlineTaskLists, err := onlineClient.ListTaskLists()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting online task lists: %v\n", err)
			os.Exit(1)
		}

		// 3. Delete all online task lists
		fmt.Println("Deleting online task lists...")
		for _, list := range onlineTaskLists.Items {
			if err := onlineClient.DeleteTaskList(gtasks.DeleteTaskListOptions{TaskListID: list.Id}); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting task list %s: %v\n", list.Id, err)
			}
		}

		// 4. Create all offline task lists online
		fmt.Println("Creating task lists...")
		for _, list := range offlineTaskLists.Items {
			createdList, err := onlineClient.CreateTaskList(gtasks.CreateTaskListOptions{Title: list.Title})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating task list %s: %v\n", list.Title, err)
				continue
			}

			offlineTasks, err := offlineClient.ListTasks(gtasks.ListTasksOptions{TaskListID: list.Id})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting offline tasks for list %s: %v\n", list.Id, err)
				continue
			}

			fmt.Printf("Creating tasks for list %s...\n", createdList.Title)
			for _, task := range offlineTasks.Items {
				_, err := onlineClient.CreateTask(gtasks.CreateTaskOptions{
					TaskListID: createdList.Id,
					Title:      task.Title,
					Notes:      task.Notes,
					Due:        task.Due,
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error creating task %s: %v\n", task.Title, err)
				}
			}
		}

		fmt.Println("Sync complete.")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}