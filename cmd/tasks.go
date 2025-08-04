package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/gtasks"
	"github.com/yanicksenn/gtasks/internal/ui"
)

var tasksCmd = &cobra.Command{
	Use:     "tasks",
	Short:   "Manage your tasks",
	Aliases: []string{"t"},
}

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks in a task list",
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		showCompleted, _ := cmd.Flags().GetBool("show-completed")
		showHidden, _ := cmd.Flags().GetBool("show-hidden")

		opts := gtasks.ListTasksOptions{
			TaskListID:    tasklist,
			ShowCompleted: showCompleted,
			ShowHidden:    showHidden,
		}

		tasks, err := h.Client.ListTasks(opts)
		if err != nil {
			return fmt.Errorf("error listing tasks: %w", err)
		}

		if !h.Quiet {
			if len(tasks.Items) == 0 {
				cmd.Println("No tasks found.")
				return nil
			}

			cmd.Println("Tasks:")
			for _, item := range tasks.Items {
				status := " "
				if item.Status == "completed" {
					status = "x"
				}
				cmd.Printf("[%s] %s (%s)\n", status, item.Title, item.Id)
			}
		}
		return nil
	},
}

var getTaskCmd = &cobra.Command{
	Use:   "get [ID]",
	Short: "Get details for a specific task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.GetTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		task, err := h.Client.GetTask(opts)
		if err != nil {
			return fmt.Errorf("error getting task: %w", err)
		}

		if !h.Quiet {
			cmd.Printf("ID:      %s\n", task.Id)
			cmd.Printf("Title:   %s\n", task.Title)
			cmd.Printf("Status:  %s\n", task.Status)
			cmd.Printf("Notes:   %s\n", task.Notes)
			cmd.Printf("Due:     %s\n", task.Due)
			cmd.Printf("Self:    %s\n", task.SelfLink)
		}
		return nil
	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	RunE: func(cmd *cobra.Command, args []string) error {
		return CreateTask(cmd, args, cmd.OutOrStdout())
	},
}

func CreateTask(cmd *cobra.Command, args []string, out io.Writer) error {
	h, err := NewCommandHelper(cmd)
	if err != nil {
		return err
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

	createdTask, err := h.Client.CreateTask(opts)
	if err != nil {
		return fmt.Errorf("error creating task: %w", err)
	}

	if !h.Quiet {
		fmt.Fprintf(out, "Successfully created task: %s (%s)\n", createdTask.Title, createdTask.Id)
	}
	return nil
}

var updateTaskCmd = &cobra.Command{
	Use:   "update [ID]",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
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

		updatedTask, err := h.Client.UpdateTask(opts)
		if err != nil {
			return fmt.Errorf("error updating task: %w", err)
		}

		if !h.Quiet {
			fmt.Printf("Successfully updated task: %s (%s)\n", updatedTask.Title, updatedTask.Id)
		}
		return nil
	},
}

var completeTaskCmd = &cobra.Command{
	Use:   "complete [ID]",
	Short: "Mark a task as complete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.CompleteTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		completedTask, err := h.Client.CompleteTask(opts)
		if err != nil {
			return fmt.Errorf("error completing task: %w", err)
		}

		if !h.Quiet {
			fmt.Printf("Successfully completed task: %s (%s)\n", completedTask.Title, completedTask.Id)
		}
		return nil
	},
}

var uncompleteTaskCmd = &cobra.Command{
	Use:   "uncomplete [ID]",
	Short: "Mark a task as not complete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.UncompleteTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		uncompletedTask, err := h.Client.UncompleteTask(opts)
		if err != nil {
			return fmt.Errorf("error uncompleting task: %w", err)
		}

		if !h.Quiet {
			fmt.Printf("Successfully uncompleted task: %s (%s)\n", uncompletedTask.Title, uncompletedTask.Id)
		}
		return nil
	},
}

var deleteTaskCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.DeleteTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		if err := h.Client.DeleteTask(opts); err != nil {
			return fmt.Errorf("error deleting task: %w", err)
		}

		if !h.Quiet {
			fmt.Printf("Successfully deleted task: %s\n", args[0])
		}
		return nil
	},
}

var printTaskCmd = &cobra.Command{
	Use:   "print [ID]",
	Short: "Print a property of a task",
	Long: `Print a property of a task.
Available properties: id, title, notes, due, status, selfLink`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.GetTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		task, err := h.Client.GetTask(opts)
		if err != nil {
			return fmt.Errorf("error getting task: %w", err)
		}

		property, _ := cmd.Flags().GetString("property")
		if err := ui.PrintTaskProperty(task, property, h.Quiet); err != nil {
			return fmt.Errorf("error printing property: %w", err)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(tasksCmd)
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