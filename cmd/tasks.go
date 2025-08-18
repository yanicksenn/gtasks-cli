package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yanicksenn/gtasks/internal/gtasks"
)

var tasksCmd = &cobra.Command{
	Use:     "tasks",
	Short:   "Manage your tasks",
	Aliases: []string{"t"},
}

func getListTasksOptions(cmd *cobra.Command) (gtasks.ListTasksOptions, gtasks.FilterOptions) {
	tasklist, _ := cmd.Flags().GetString("tasklist")
	showCompleted, _ := cmd.Flags().GetBool("show-completed")
	showHidden, _ := cmd.Flags().GetBool("show-hidden")
	titleContains, _ := cmd.Flags().GetString("title-contains")
	notesContains, _ := cmd.Flags().GetString("notes-contains")
	dueBefore, _ := cmd.Flags().GetString("due-before")
	dueAfter, _ := cmd.Flags().GetString("due-after")
	sortBy, _ := cmd.Flags().GetString("sort-by")

	listOpts := gtasks.ListTasksOptions{
		TaskListID:    tasklist,
		ShowCompleted: showCompleted,
		ShowHidden:    showHidden,
		SortBy:        sortBy,
	}

	filterOpts := gtasks.FilterOptions{
		TitleContains: titleContains,
		NotesContains: notesContains,
		DueBefore:     dueBefore,
		DueAfter:      dueAfter,
	}

	return listOpts, filterOpts
}

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks in a task list",
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		listOpts, filterOpts := getListTasksOptions(cmd)

		tasks, err := h.Client.ListTasks(listOpts)
		if err != nil {
			return fmt.Errorf("error listing tasks: %w", err)
		}

		filteredTasks, err := gtasks.FilterTasks(tasks.Items, filterOpts)
		if err != nil {
			return fmt.Errorf("error filtering tasks: %w", err)
		}
		tasks.Items = filteredTasks

		return h.Printer.PrintTasks(tasks)
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

		// Get the tasklist ID from the flags
		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.GetTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		// Get the task
		task, err := h.Client.GetTask(opts)
		if err != nil {
			return fmt.Errorf("error getting task: %w", err)
		}

		// Print the task
		return h.Printer.PrintTask(task)
	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := NewCommandHelper(cmd)
		if err != nil {
			return err
		}

		// Get the flag values
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

		// Create the task
		createdTask, err := h.Client.CreateTask(opts)
		if err != nil {
			return fmt.Errorf("error creating task: %w", err)
		}

		// Print the created task
		return h.Printer.PrintTask(createdTask)
	},
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

		// Get the flag values
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

		// Update the task
		updatedTask, err := h.Client.UpdateTask(opts)
		if err != nil {
			return fmt.Errorf("error updating task: %w", err)
		}

		// Print a success message
		return h.Printer.PrintSuccess(fmt.Sprintf("Successfully updated task: %s (%s)", updatedTask.Title, updatedTask.Id))
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

		// Get the tasklist ID from the flags
		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.CompleteTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		// Complete the task
		completedTask, err := h.Client.CompleteTask(opts)
		if err != nil {
			return fmt.Errorf("error completing task: %w", err)
		}

		// Print a success message
		return h.Printer.PrintSuccess(fmt.Sprintf("Successfully completed task: %s (%s)", completedTask.Title, completedTask.Id))
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

		// Get the tasklist ID from the flags
		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.UncompleteTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		// Uncomplete the task
		uncompletedTask, err := h.Client.UncompleteTask(opts)
		if err != nil {
			return fmt.Errorf("error uncompleting task: %w", err)
		}

		// Print a success message
		return h.Printer.PrintSuccess(fmt.Sprintf("Successfully uncompleted task: %s (%s)", uncompletedTask.Title, uncompletedTask.Id))
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

		// Get the tasklist ID from the flags
		tasklist, _ := cmd.Flags().GetString("tasklist")
		opts := gtasks.DeleteTaskOptions{
			TaskListID: tasklist,
			TaskID:     args[0],
		}

		// Delete the task
		if err := h.Client.DeleteTask(opts); err != nil {
			return fmt.Errorf("error deleting task: %w", err)
		}

		// Print a success message
		return h.Printer.PrintDelete("task", args[0])
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

	listTasksCmd.Flags().String("tasklist", "@default", "The ID of the task list")
	listTasksCmd.Flags().Bool("show-completed", false, "Include completed tasks in the output")
	listTasksCmd.Flags().Bool("show-hidden", false, "Include hidden tasks in the output")
	listTasksCmd.Flags().String("title-contains", "", "Filter tasks by title (case-insensitive)")
	listTasksCmd.Flags().String("notes-contains", "", "Filter tasks by notes (case-insensitive)")
	listTasksCmd.Flags().String("due-before", "", "Filter tasks with a due date before the specified date (e.g., '2025-12-31')")
	listTasksCmd.Flags().String("due-after", "", "Filter tasks with a due date after the specified date (e.g., '2025-12-31')")
	listTasksCmd.Flags().String("sort-by", "alphabetical", "Sort tasks by (alphabetical, last-modified, due-date)")

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
}