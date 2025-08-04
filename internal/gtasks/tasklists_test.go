package gtasks

import (
	"strings"
	"testing"

	"github.com/yanicksenn/gtasks/internal/store"
	"github.com/yanicksenn/gtasks/internal/ui"
)

func newTestOfflineClient(t *testing.T) *offlineClient {
	t.Helper()
	return &offlineClient{
		store: store.NewTestStore(),
	}
}

func TestTaskListsLifecycle(t *testing.T) {
	client := newTestOfflineClient(t)

	// 1. Initial list should contain the default list
	output := CaptureOutput(t, func() {
		lists, err := client.ListTaskLists(ListTaskListsOptions{})
		if err != nil {
			t.Fatalf("ListTaskLists failed: %v", err)
		}
		printTaskLists(lists)
	})
	if !strings.Contains(output, "Default List") {
		t.Errorf("expected default list, got '%s'", output)
	}

	// 2. Create a new task list
	var listID string
	CaptureOutput(t, func() {
		opts := CreateTaskListOptions{Title: "Groceries"}
		list, err := client.CreateTaskList(opts)
		if err != nil {
			t.Fatalf("CreateTaskList failed: %v", err)
		}
		listID = list.Id
	})

	// 3. List should now contain the new list
	output = CaptureOutput(t, func() {
		lists, err := client.ListTaskLists(ListTaskListsOptions{})
		if err != nil {
			t.Fatalf("ListTaskLists failed: %v", err)
		}
		printTaskLists(lists)
	})
	if !strings.Contains(output, "Groceries") {
		t.Errorf("expected list to contain 'Groceries', got '%s'", output)
	}

	// 4. Get the list by ID
	output = CaptureOutput(t, func() {
		opts := GetTaskListOptions{TaskListID: listID}
		list, err := client.GetTaskList(opts)
		if err != nil {
			t.Fatalf("GetTaskList failed: %v", err)
		}
		printTaskList(list)
	})
	if !strings.Contains(output, "Title: Groceries") {
		t.Errorf("expected get to show 'Title: Groceries', got '%s'", output)
	}

	// 5. Update the list
	CaptureOutput(t, func() {
		opts := UpdateTaskListOptions{TaskListID: listID, Title: "Updated Groceries"}
		_, err := client.UpdateTaskList(opts)
		if err != nil {
			t.Fatalf("UpdateTaskList failed: %v", err)
		}
	})

	// 6. Delete the list
	CaptureOutput(t, func() {
		opts := DeleteTaskListOptions{TaskListID: listID}
		err := client.DeleteTaskList(opts)
		if err != nil {
			t.Fatalf("DeleteTaskList failed: %v", err)
		}
	})

	// 7. Final list should not contain the deleted list
	output = CaptureOutput(t, func() {
		lists, err := client.ListTaskLists(ListTaskListsOptions{})
		if err != nil {
			t.Fatalf("ListTaskLists failed: %v", err)
		}
		printTaskLists(lists)
	})
	if strings.Contains(output, "Groceries") {
		t.Errorf("expected list to not contain 'Groceries', got '%s'", output)
	}
}

func TestTaskListPrint(t *testing.T) {
	client := newTestOfflineClient(t)

	// 1. Create a new task list
	var listID string
	CaptureOutput(t, func() {
		opts := CreateTaskListOptions{Title: "Test List"}
		list, err := client.CreateTaskList(opts)
		if err != nil {
			t.Fatalf("CreateTaskList failed: %v", err)
		}
		listID = list.Id
	})

	// 2. Print the title
	output := CaptureOutput(t, func() {
		list, err := client.GetTaskList(GetTaskListOptions{TaskListID: listID})
		if err != nil {
			t.Fatalf("GetTaskList failed: %v", err)
		}
		if err := ui.PrintTaskListProperty(list, "title", false); err != nil {
			t.Fatalf("PrintTaskListProperty failed: %v", err)
		}
	})
	if !strings.Contains(output, "Test List") {
		t.Errorf("expected output to contain 'Test List', got '%s'", output)
	}
}
