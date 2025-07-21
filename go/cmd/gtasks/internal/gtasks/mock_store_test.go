package gtasks

import (
	"fmt"
	"sync"

	"google.golang.org/api/tasks/v1"
)

// mockStore is a thread-safe, in-memory database for the mock server.
type mockStore struct {
	mu        sync.Mutex
	taskLists map[string]*tasks.TaskList
	tasks     map[string]map[string]*tasks.Task // taskListID -> taskID -> task
	nextID    int
}

func newMockStore() *mockStore {
	s := &mockStore{
		taskLists: make(map[string]*tasks.TaskList),
		tasks:     make(map[string]map[string]*tasks.Task),
		nextID:    1,
	}
	// Pre-populate with a default list for testing task creation without first creating a list
	s.createTaskList(&tasks.TaskList{Title: "Default List"})
	return s
}

func (s *mockStore) newID() string {
	id := fmt.Sprintf("id%d", s.nextID)
	s.nextID++
	return id
}

func (s *mockStore) createTaskList(list *tasks.TaskList) *tasks.TaskList {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	newList := &tasks.TaskList{
		Id:    id,
		Title: list.Title,
	}
	s.taskLists[id] = newList
	s.tasks[id] = make(map[string]*tasks.Task)
	return newList
}

func (s *mockStore) getTaskList(id string) *tasks.TaskList {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.taskLists[id]
}

func (s *mockStore) listTaskLists() []*tasks.TaskList {
	s.mu.Lock()
	defer s.mu.Unlock()
	var lists []*tasks.TaskList
	for _, list := range s.taskLists {
		lists = append(lists, list)
	}
	return lists
}

func (s *mockStore) updateTaskList(id string, list *tasks.TaskList) *tasks.TaskList {
	s.mu.Lock()
	defer s.mu.Unlock()
	existingList := s.taskLists[id]
	existingList.Title = list.Title
	return existingList
}

func (s *mockStore) deleteTaskList(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.taskLists, id)
	delete(s.tasks, id)
}

func (s *mockStore) createTask(listID string, task *tasks.Task) *tasks.Task {
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
	if _, ok := s.tasks[listID]; !ok {
		s.tasks[listID] = make(map[string]*tasks.Task)
	}
	s.tasks[listID][id] = newTask
	return newTask
}

func (s *mockStore) getTask(listID, taskID string) *tasks.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tasks[listID][taskID]
}

func (s *mockStore) listTasks(listID string) []*tasks.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	var tasks []*tasks.Task
	for _, task := range s.tasks[listID] {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *mockStore) updateTask(listID, taskID string, task *tasks.Task) *tasks.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	existingTask := s.tasks[listID][taskID]
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
	return existingTask
}

func (s *mockStore) deleteTask(listID, taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks[listID], taskID)
}