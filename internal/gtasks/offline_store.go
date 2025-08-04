package gtasks

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

// offlineStore manages the state of tasks and task lists in a local file.
type offlineStore struct {
	mu   sync.Mutex
	path string
	Data struct {
		TaskLists map[string]*tasks.TaskList       `json:"task_lists"`
		Tasks     map[string]map[string]*tasks.Task `json:"tasks"` // taskListID -> taskID -> task
		NextID    int                              `json:"next_id"`
	} `json:"data"`
}

// newOfflineStore creates a new offline store, loading data from the file if it exists.
func newOfflineStore() (*offlineStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".config", "gtasks", offlineDataFile)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	store := &offlineStore{path: path}
	store.Data.TaskLists = make(map[string]*tasks.TaskList)
	store.Data.Tasks = make(map[string]map[string]*tasks.Task)
	store.Data.NextID = 1

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

func (s *offlineStore) persist() error {
	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

func (s *offlineStore) newID() string {
	id := fmt.Sprintf("offline-id%d", s.Data.NextID)
	s.Data.NextID++
	return id
}

func (s *offlineStore) createTaskList(list *tasks.TaskList) (*tasks.TaskList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	newList := &tasks.TaskList{Id: id, Title: list.Title}
	s.Data.TaskLists[id] = newList
	s.Data.Tasks[id] = make(map[string]*tasks.Task)

	return newList, s.persist()
}

func (s *offlineStore) listTaskLists() ([]*tasks.TaskList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var lists []*tasks.TaskList
	for _, list := range s.Data.TaskLists {
		lists = append(lists, list)
	}
	return lists, nil
}

// ... (other methods for tasks and tasklists to be implemented) ...

func (s *offlineStore) getTaskList(id string) (*tasks.TaskList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Data.TaskLists[id], nil
}

func (s *offlineStore) updateTaskList(id string, list *tasks.TaskList) (*tasks.TaskList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existingList := s.Data.TaskLists[id]
	existingList.Title = list.Title
	return existingList, s.persist()
}

func (s *offlineStore) deleteTaskList(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Data.TaskLists, id)
	delete(s.Data.Tasks, id)
	return s.persist()
}

func (s *offlineStore) createTask(listID string, task *tasks.Task) (*tasks.Task, error) {
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

func (s *offlineStore) getTask(listID, taskID string) (*tasks.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Data.Tasks[listID][taskID], nil
}

func (s *offlineStore) listTasks(listID string) ([]*tasks.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var tasks []*tasks.Task
	for _, task := range s.Data.Tasks[listID] {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *offlineStore) updateTask(listID, taskID string, task *tasks.Task) (*tasks.Task, error) {
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

func (s *offlineStore) deleteTask(listID, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Data.Tasks[listID], taskID)
	return s.persist()
}
