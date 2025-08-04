package gtasks

import (
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/api/tasks/v1"
)

func TestOfflineStore(t *testing.T) {
	// Create a temporary directory for the test.
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// 1. Create a new store
	store, err := newOfflineStore()
	if err != nil {
		t.Fatalf("failed to create offline store: %v", err)
	}

	// Verify the file was created.
	expectedPath := filepath.Join(tempDir, ".config", "gtasks", offlineDataFile)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("offline data file was not created at %s", expectedPath)
	}

	// 2. Create a task list
	list, err := store.createTaskList(&tasks.TaskList{Title: "Test List"})
	if err != nil {
		t.Fatalf("failed to create task list: %v", err)
	}
	if list.Title != "Test List" {
		t.Errorf("expected title 'Test List', got '%s'", list.Title)
	}

	// 3. Create a task
	task, err := store.createTask(list.Id, &tasks.Task{Title: "Test Task"})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	if task.Title != "Test Task" {
		t.Errorf("expected title 'Test Task', got '%s'", task.Title)
	}

	// 4. Re-load the store from the file to verify persistence
	store2, err := newOfflineStore()
	if err != nil {
		t.Fatalf("failed to re-load offline store: %v", err)
	}

	// 5. Verify the re-loaded data
	reloadedList, err := store2.getTaskList(list.Id)
	if err != nil {
		t.Fatalf("failed to get re-loaded task list: %v", err)
	}
	if reloadedList.Title != "Test List" {
		t.Errorf("expected re-loaded title 'Test List', got '%s'", reloadedList.Title)
	}

	reloadedTask, err := store2.getTask(list.Id, task.Id)
	if err != nil {
		t.Fatalf("failed to get re-loaded task: %v", err)
	}
	if reloadedTask.Title != "Test Task" {
		t.Errorf("expected re-loaded title 'Test Task', got '%s'", reloadedTask.Title)
	}
}
