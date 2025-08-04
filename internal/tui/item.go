package tui

import (
	taskspb "google.golang.org/api/tasks/v1"
)

type taskListItem struct {
	*taskspb.TaskList
}

func (i taskListItem) FilterValue() string {
	return i.TaskList.Title
}

func (i taskListItem) Title() string {
	return i.TaskList.Title
}

func (i taskListItem) Description() string {
	return i.TaskList.Id
}

type taskItem struct {
	*taskspb.Task
}

func (i taskItem) FilterValue() string {
	return i.Task.Title
}

func (i taskItem) Title() string {
	return i.Task.Title
}

func (i taskItem) Description() string {
	return i.Task.Notes
}

