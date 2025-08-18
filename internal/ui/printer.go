package ui

import (
	"encoding/json"
	"fmt"
	"io"

	"google.golang.org/api/tasks/v1"
	"gopkg.in/yaml.v3"
)

// OutputFormat defines the output format for the printer.
type OutputFormat string

const (
	// JSONFormat is the JSON output format.
	JSONFormat OutputFormat = "json"
	// YAMLFormat is the YAML output format.
	YAMLFormat OutputFormat = "yaml"
	// TableFormat is the table output format.
	TableFormat OutputFormat = "table"
)

// Printer handles formatting and printing data to the console.
type Printer struct {
	out    io.Writer
	format OutputFormat
	quiet  bool
}

// NewPrinter creates a new Printer.
func NewPrinter(out io.Writer, format string, quiet bool) *Printer {
	var f OutputFormat
	switch format {
	case "json":
		f = JSONFormat
	case "yaml":
		f = YAMLFormat
	default:
		f = TableFormat
	}
	return &Printer{out: out, format: f, quiet: quiet}
}

// PrintTaskLists prints a list of task lists.
func (p *Printer) PrintTaskLists(lists *tasks.TaskLists) error {
	switch p.format {
	case JSONFormat:
		return json.NewEncoder(p.out).Encode(lists)
	case YAMLFormat:
		return yaml.NewEncoder(p.out).Encode(lists)
	default:
		if p.quiet {
			return nil
		}
		if len(lists.Items) == 0 {
			fmt.Fprintln(p.out, "No task lists found.")
			return nil
		}
		fmt.Fprintln(p.out, "Task Lists:")
		for _, item := range lists.Items {
			fmt.Fprintf(p.out, "- %s (%s)\n", item.Title, item.Id)
		}
		return nil
	}
}

// PrintTaskList prints a single task list.
func (p *Printer) PrintTaskList(list *tasks.TaskList) error {
	switch p.format {
	case JSONFormat:
		return json.NewEncoder(p.out).Encode(list)
	case YAMLFormat:
		return yaml.NewEncoder(p.out).Encode(list)
	default:
		if p.quiet {
			return nil
		}
		fmt.Fprintf(p.out, "ID:    %s\n", list.Id)
		fmt.Fprintf(p.out, "Title: %s\n", list.Title)
		fmt.Fprintf(p.out, "Self:  %s\n", list.SelfLink)
		return nil
	}
}

// PrintTasks prints a list of tasks.
func (p *Printer) PrintTasks(tasks *tasks.Tasks) error {
	switch p.format {
	case JSONFormat:
		return json.NewEncoder(p.out).Encode(tasks)
	case YAMLFormat:
		return yaml.NewEncoder(p.out).Encode(tasks)
	default:
		if p.quiet {
			return nil
		}
		if len(tasks.Items) == 0 {
			fmt.Fprintln(p.out, "No tasks found.")
			return nil
		}
		fmt.Fprintln(p.out, "Tasks:")
		for _, item := range tasks.Items {
			status := " "
			if item.Status == "completed" {
				status = "x"
			}
			fmt.Fprintf(p.out, "[%s] %s (%s)\n", status, item.Title, item.Id)
		}
		return nil
	}
}

// PrintTask prints a single task.
func (p *Printer) PrintTask(task *tasks.Task) error {
	switch p.format {
	case JSONFormat:
		return json.NewEncoder(p.out).Encode(task)
	case YAMLFormat:
		return yaml.NewEncoder(p.out).Encode(task)
	default:
		if p.quiet {
			return nil
		}
		fmt.Fprintf(p.out, "ID:      %s\n", task.Id)
		fmt.Fprintf(p.out, "Title:   %s\n", task.Title)
		fmt.Fprintf(p.out, "Status:  %s\n", task.Status)
		fmt.Fprintf(p.out, "Notes:   %s\n", task.Notes)
		fmt.Fprintf(p.out, "Due:     %s\n", task.Due)
		fmt.Fprintf(p.out, "Self:    %s\n", task.SelfLink)
		return nil
	}
}

// PrintAccounts prints a list of accounts.
func (p *Printer) PrintAccounts(accounts []string, activeAccount string) error {
	switch p.format {
	case JSONFormat:
		return json.NewEncoder(p.out).Encode(accounts)
	case YAMLFormat:
		return yaml.NewEncoder(p.out).Encode(accounts)
	default:
		if p.quiet {
			return nil
		}
		if len(accounts) == 0 {
			fmt.Fprintln(p.out, "No accounts authenticated.")
			return nil
		}
		fmt.Fprintln(p.out, "Authenticated Accounts:")
		for _, account := range accounts {
			if account == activeAccount {
				fmt.Fprintf(p.out, "- %s (active)\n", account)
			} else {
				fmt.Fprintf(p.out, "- %s\n", account)
			}
		}
		return nil
	}
}

// PrintSuccess prints a success message.
func (p *Printer) PrintSuccess(message string) error {
	if p.quiet {
		return nil
	}
	fmt.Fprintln(p.out, message)
	return nil
}

// PrintDelete prints a delete message.
func (p *Printer) PrintDelete(resource string, id string) error {
	if p.quiet {
		return nil
	}
	fmt.Fprintf(p.out, "Successfully deleted %s: %s\n", resource, id)
	return nil
}

