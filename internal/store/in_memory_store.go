package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/api/tasks/v1"
)

const (
	offlineDataFile = "offline.json"
)

// InMemoryStore manages the state of tasks and task lists in memory,
// with an option to persist to a local file.
type InMemoryStore struct {
	mu   sync.Mutex
	path string // Path for persistence; if empty, store is transient.
	Data struct {
		TaskLists map[string]*tasks.TaskList       `json:"task_lists"`
		Tasks     map[string]map[string]*tasks.Task `json:"tasks"` // taskListID -> taskID -> task
		NextID    int                              `json:"next_id"`
	} `json:"data"`
}

// NewInMemoryStore creates a new in-memory store. If a path is provided,
// it loads data from that file if it exists.
func NewInMemoryStore(path string) (*InMemoryStore, error) {
	store := &InMemoryStore{path: path}
	store.Data.TaskLists = make(map[string]*tasks.TaskList)
	store.Data.Tasks = make(map[string]map[string]*tasks.Task)
	store.Data.NextID = 1

	if path == "" {
		return store, nil // Transient store
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return store, store.persist() // Create the file if it doesn't exist
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 { // File is empty
		return store, store.persist()
	}

	if err := json.Unmarshal(data, &store.Data); err != nil {
		return nil, err
	}

	return store, nil
}

// NewTestStore creates a new, empty, transient in-memory store for testing.
func NewTestStore() *InMemoryStore {
	store, _ := NewInMemoryStore("")
	// Pre-populate with a default list for testing
	store.CreateTaskList(&tasks.TaskList{Title: "Default List"})
	return store
}

// GetOfflineStorePath returns the default path for the offline data file.
func GetOfflineStorePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gtasks", offlineDataFile), nil
}

// persist writes the current state of the store to the file system.
func (s *InMemoryStore) persist() error {
	if s.path == "" {
		return nil // Don't persist for a transient store
	}
	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

// newID generates a new unique ID for a task or task list.
func (s *InMemoryStore) newID() string {
	id := fmt.Sprintf("id%d", s.Data.NextID)
	s.Data.NextID++
	return id
}

// CreateTaskList creates a new task list in the store.
func (s *InMemoryStore) CreateTaskList(list *tasks.TaskList) (*tasks.TaskList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	newList := &tasks.TaskList{Id: id, Title: list.Title}
	s.Data.TaskLists[id] = newList
	s.Data.Tasks[id] = make(map[string]*tasks.Task)

	return newList, s.persist()
}

// ListTaskLists returns all the task lists in the store.
func (s *InMemoryStore) ListTaskLists() ([]*tasks.TaskList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var lists []*tasks.TaskList
	for _, list := range s.Data.TaskLists {
		lists = append(lists, list)
	}
	return lists, nil
}

// GetTaskList returns a task list from the store by its ID.
func (s *InMemoryStore) GetTaskList(id string) (*tasks.TaskList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Data.TaskLists[id], nil
}

// UpdateTaskList updates a task list in the store.
func (s *InMemoryStore) UpdateTaskList(id string, list *tasks.TaskList) (*tasks.TaskList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existingList := s.Data.TaskLists[id]
	existingList.Title = list.Title
	return existingList, s.persist()
}

// DeleteTaskList deletes a task list from the store.
func (s *InMemoryStore) DeleteTaskList(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Data.TaskLists, id)
	delete(s.Data.Tasks, id)
	return s.persist()
}

// CreateTask creates a new task in the store.
func (s *InMemoryStore) CreateTask(listID string, task *tasks.Task) (*tasks.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	newTask := &tasks.Task{
		Id:     id,
		Title:  task.Title,
		Notes:  task.Notes,
		Due:    task.Due,
		Status: "needsAction",
	}
	if _, ok := s.Data.Tasks[listID]; !ok {
		s.Data.Tasks[listID] = make(map[string]*tasks.Task)
	}
	s.Data.Tasks[listID][id] = newTask
	return newTask, s.persist()
}

// GetTask returns a task from the store by its ID.
func (s *InMemoryStore) GetTask(listID, taskID string) (*tasks.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Data.Tasks[listID][taskID], nil
}

// ListTasks returns all the tasks in a given task list.
func (s *InMemoryStore) ListTasks(listID string) ([]*tasks.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var tasks []*tasks.Task
	for _, task := range s.Data.Tasks[listID] {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// UpdateTask updates a task in the store.
func (s *InMemoryStore) UpdateTask(listID, taskID string, task *tasks.Task) (*tasks.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existingTask := s.Data.Tasks[listID][taskID]
	if task.Title != "" {
		existingTask.Title = task.Title
	}
	if task.Notes != "" {
		existingTask.Notes = task.Notes
	}
	if task.Due != "" {
		existingTask.Due = task.Due
	}
	if task.Status != "" {
		existingTask.Status = task.Status
	}
	return existingTask, s.persist()
}

// DeleteTask deletes a task from the store.
func (s *InMemoryStore) DeleteTask(listID, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Data.Tasks[listID], taskID)
	return s.persist()
}
