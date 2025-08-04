package gtasks

import (
	"strings"
	"time"

	taskspb "google.golang.org/api/tasks/v1"
)

// FilterOptions holds the criteria for filtering tasks.
type FilterOptions struct {
	TitleContains string
	NotesContains string
	DueBefore     string
	DueAfter      string
}

// FilterTasks filters a slice of tasks based on the provided options.
func FilterTasks(tasks []*taskspb.Task, opts FilterOptions) ([]*taskspb.Task, error) {
	var filtered []*taskspb.Task

	for _, task := range tasks {
		if opts.TitleContains != "" && !strings.Contains(strings.ToLower(task.Title), strings.ToLower(opts.TitleContains)) {
			continue
		}

		if opts.NotesContains != "" && !strings.Contains(strings.ToLower(task.Notes), strings.ToLower(opts.NotesContains)) {
			continue
		}

		if opts.DueBefore != "" {
			dueBefore, err := time.Parse("2006-01-02", opts.DueBefore)
			if err != nil {
				return nil, err
			}
			taskDue, err := time.Parse(time.RFC3339, task.Due)
			if err != nil {
				continue
			}
			if !taskDue.Before(dueBefore) {
				continue
			}
		}

		if opts.DueAfter != "" {
			dueAfter, err := time.Parse("2006-01-02", opts.DueAfter)
			if err != nil {
				return nil, err
			}
			taskDue, err := time.Parse(time.RFC3339, task.Due)
			if err != nil {
				continue
			}
			if !taskDue.After(dueAfter) {
				continue
			}
		}

		filtered = append(filtered, task)
	}

	return filtered, nil
}