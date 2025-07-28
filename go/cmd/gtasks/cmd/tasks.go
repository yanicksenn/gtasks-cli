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

		tasks, err := client.ListTasks(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing tasks: %v\n", err)
			os.Exit(1)
		}

		if len(tasks.Items) == 0 {
			fmt.Println("No tasks found.")
			return
		}

		fmt.Println("Tasks:")
		for _, item := range tasks.Items {
			status := " "
			if item.Status == "completed" {
				status = "x"
			}
			fmt.Printf("[%s] %s (%s)\n", status, item.Title, item.Id)
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

		task, err := client.GetTask(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("ID:      %s\n", task.Id)
		fmt.Printf("Title:   %s\n", task.Title)
		fmt.Printf("Status:  %s\n", task.Status)
		fmt.Printf("Notes:   %s\n", task.Notes)
		fmt.Printf("Due:     %s\n", task.Due)
		fmt.Printf("Self:    %s\n", task.SelfLink)
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

		createdTask, err := client.CreateTask(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created task: %s (%s)\n", createdTask.Title, createdTask.Id)
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

		updatedTask, err := client.UpdateTask(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully updated task: %s (%s)\n", updatedTask.Title, updatedTask.Id)
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

		completedTask, err := client.CompleteTask(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error completing task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully completed task: %s (%s)\n", completedTask.Title, completedTask.Id)
	},
}

var uncompleteTaskCmd = &cobra.Command{
	Use:   "uncomplete [ID]",
	Short: "Mark a task as not complete",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(cmd, context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
			os.Exit(1)
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.UncompleteTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		uncompletedTask, err := client.UncompleteTask(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error uncompleting task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully uncompleted task: %s (%s)\n", uncompletedTask.Title, uncompletedTask.Id)
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

		fmt.Printf("Successfully deleted task: %s\n", args[0])
	},
}

var printTaskCmd = &cobra.Command{
	Use:   "print [ID]",
	Short: "Print a property of a task",
	Long: `Print a property of a task.
Available properties: id, title, notes, due, status, selfLink`,
	Args: cobra.ExactArgs(1),
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

		task, err := client.GetTask(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting task: %v\n", err)
			os.Exit(1)
		}

		property, _ := cmd.Flags().GetString("property")
		if err := gtasks.PrintTaskProperty(task, property); err != nil {
			fmt.Fprintf(os.Stderr, "Error printing property: %v\n", err)
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
	tasksCmd.AddCommand(uncompleteTaskCmd)
	tasksCmd.AddCommand(deleteTaskCmd)
	tasksCmd.AddCommand(printTaskCmd)

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

	uncompleteTaskCmd.Flags().String("tasklist", "@default", "The ID of the task list")

	deleteTaskCmd.Flags().String("tasklist", "@default", "The ID of the task list")

	printTaskCmd.Flags().String("tasklist", "@default", "The ID of the task list")
	printTaskCmd.Flags().String("property", "", "The property to print")
	printTaskCmd.MarkFlagRequired("property")
}