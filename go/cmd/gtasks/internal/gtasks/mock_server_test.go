package gtasks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// newMockServer creates a new mock server that simulates the Google Tasks API.
func newMockServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock for task lists
	mux.HandleFunc("/tasks/v1/users/@me/lists", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"kind": "tasks#taskLists",
				"items": [
					{
						"kind": "tasks#taskList",
						"id": "taskList1",
						"title": "Test Task List 1"
					},
					{
						"kind": "tasks#taskList",
						"id": "taskList2",
						"title": "Test Task List 2"
					}
				]
			}`)
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"kind": "tasks#taskList",
				"id": "newTaskList",
				"title": "New Task List"
			}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Mock for a single task list
	mux.HandleFunc("/tasks/v1/users/@me/lists/taskList1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"kind": "tasks#taskList",
				"id": "taskList1",
				"title": "Test Task List 1"
			}`)
		case http.MethodPut:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"kind": "tasks#taskList",
				"id": "taskList1",
				"title": "Updated Task List"
			}`)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Mock for tasks in a task list
	mux.HandleFunc("/tasks/v1/lists/taskList1/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"kind": "tasks#tasks",
				"items": [
					{
						"kind": "tasks#task",
						"id": "task1",
						"title": "Test Task 1"
					},
					{
						"kind": "tasks#task",
						"id": "task2",
						"title": "Test Task 2"
					}
				]
			}`)
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"kind": "tasks#task",
				"id": "newTask",
				"title": "New Task"
			}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Mock for a single task
	mux.HandleFunc("/tasks/v1/lists/taskList1/tasks/task1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"kind": "tasks#task",
				"id": "task1",
				"title": "Test Task 1"
			}`)
		case http.MethodPut:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{
				"kind": "tasks#task",
				"id": "task1",
				"title": "Updated Task"
			}`)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Mock for completing a task (get the task first)
	mux.HandleFunc("/tasks/v1/lists/taskList1/tasks/task1/complete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"kind": "tasks#task",
			"id": "task1",
			"title": "Test Task 1"
		}`)
	})

	return httptest.NewServer(mux)
}