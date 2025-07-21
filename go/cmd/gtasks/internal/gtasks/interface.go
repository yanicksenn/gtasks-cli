package gtasks

import "google.golang.org/api/tasks/v1"

// TasksService is an interface that wraps the Google Tasks API service.
// This is used to allow mocking in tests.
type TasksService interface {
	// Tasklists is the service for managing task lists.
	Tasklists() *tasks.TasklistsService
	// Tasks is the service for managing tasks.
	Tasks() *tasks.TasksService
}

// TasksServiceWrapper is a wrapper for the Google Tasks API service that implements the TasksService interface.
type TasksServiceWrapper struct {
	service *tasks.Service
}

// Tasklists returns the task lists service.
func (w *TasksServiceWrapper) Tasklists() *tasks.TasklistsService {
	return w.service.Tasklists
}

// Tasks returns the tasks service.
func (w *TasksServiceWrapper) Tasks() *tasks.TasksService {
	return w.service.Tasks
}
