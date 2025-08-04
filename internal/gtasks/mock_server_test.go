package gtasks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"google.golang.org/api/tasks/v1"
)

func newMockServer() *httptest.Server {
	store := newMockStore()
	mux := http.NewServeMux()

	// TaskLists
	mux.HandleFunc("/tasks/v1/users/@me/lists", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			lists := store.listTaskLists()
			json.NewEncoder(w).Encode(&tasks.TaskLists{Items: lists})
		case http.MethodPost:
			var list tasks.TaskList
			json.NewDecoder(r.Body).Decode(&list)
			createdList := store.createTaskList(&list)
			json.NewEncoder(w).Encode(createdList)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/v1/users/@me/lists/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Extract ID: /tasks/v1/users/@me/lists/{listID}
		id := strings.TrimPrefix(r.URL.Path, "/tasks/v1/users/@me/lists/")
		switch r.Method {
		case http.MethodGet:
			list := store.getTaskList(id)
			json.NewEncoder(w).Encode(list)
		case http.MethodPut:
			var list tasks.TaskList
			json.NewDecoder(r.Body).Decode(&list)
			updatedList := store.updateTaskList(id, &list)
			json.NewEncoder(w).Encode(updatedList)
		case http.MethodDelete:
			store.deleteTaskList(id)
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Tasks
	mux.HandleFunc("/tasks/v1/lists/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Extract IDs: /tasks/v1/lists/{listID}/tasks/{taskID}
		path := strings.TrimPrefix(r.URL.Path, "/tasks/v1/lists/")
		parts := strings.Split(path, "/")

		listID := parts[0]
		isTaskRequest := len(parts) > 1 && parts[1] == "tasks"

		if !isTaskRequest {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Handle list-level task operations (/tasks/v1/lists/{listID}/tasks)
		if len(parts) == 2 {
			switch r.Method {
			case http.MethodGet:
				taskItems := store.listTasks(listID)
				json.NewEncoder(w).Encode(&tasks.Tasks{Items: taskItems})
			case http.MethodPost:
				var task tasks.Task
				json.NewDecoder(r.Body).Decode(&task)
				createdTask := store.createTask(listID, &task)
				json.NewEncoder(w).Encode(createdTask)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}

		// Handle specific task operations (/tasks/v1/lists/{listID}/tasks/{taskID})
		if len(parts) == 3 {
			taskID := parts[2]
			switch r.Method {
			case http.MethodGet:
				task := store.getTask(listID, taskID)
				json.NewEncoder(w).Encode(task)
			case http.MethodPut:
				var task tasks.Task
				json.NewDecoder(r.Body).Decode(&task)
				updatedTask := store.updateTask(listID, taskID, &task)
				json.NewEncoder(w).Encode(updatedTask)
			case http.MethodDelete:
				store.deleteTask(listID, taskID)
				w.WriteHeader(http.StatusNoContent)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	})

	return httptest.NewServer(mux)
}