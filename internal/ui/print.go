package ui

import (
	"fmt"
	"strings"

	"google.golang.org/api/tasks/v1"
)

// PrintTaskListProperty prints a specific property of a task list to stdout.
func PrintTaskListProperty(list *tasks.TaskList, property string, quiet bool) error {
	if quiet {
		return nil
	}
	switch strings.ToLower(property) {
	case "id":
		fmt.Println(list.Id)
	case "title":
		fmt.Println(list.Title)
	case "selflink":
		fmt.Println(list.SelfLink)
	default:
		return fmt.Errorf("unknown property: %s. Available properties: id, title, selfLink", property)
	}
	return nil
}

// PrintTaskProperty prints a specific property of a task to stdout.
func PrintTaskProperty(task *tasks.Task, property string, quiet bool) error {
	if quiet {
		return nil
	}
	switch strings.ToLower(property) {
	case "id":
		fmt.Println(task.Id)
	case "title":
		fmt.Println(task.Title)
	case "notes":
		fmt.Println(task.Notes)
	case "due":
		fmt.Println(task.Due)
	case "status":
		fmt.Println(task.Status)
	case "selflink":
		fmt.Println(task.SelfLink)
	default:
		return fmt.Errorf("unknown property: %s. Available properties: id, title, notes, due, status, selfLink", property)
	}
	return nil
}
