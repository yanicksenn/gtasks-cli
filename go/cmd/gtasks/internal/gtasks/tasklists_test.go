package gtasks

import (
	"strings"
	"testing"
)

func TestTaskListsLifecycle(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// 1. Initial list should contain the default list
	output := captureOutput(t, func() {
		lists, err := client.ListTaskLists()
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
	output = captureOutput(t, func() {
		opts := CreateTaskListOptions{Title: "Groceries"}
		list, err := client.CreateTaskList(opts)
		if err != nil {
			t.Fatalf("CreateTaskList failed: %v", err)
		}
		listID = list.Id
		printTaskList(list)
	})

	// 3. List should now contain the new list
	output = captureOutput(t, func() {
		lists, err := client.ListTaskLists()
		if err != nil {
			t.Fatalf("ListTaskLists failed: %v", err)
		}
		printTaskLists(lists)
	})
	if !strings.Contains(output, "Groceries") {
		t.Errorf("expected list to contain 'Groceries', got '%s'", output)
	}

	// 4. Get the list by ID
	output = captureOutput(t, func() {
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
	output = captureOutput(t, func() {
		opts := UpdateTaskListOptions{TaskListID: listID, Title: "Updated Groceries"}
		_, err := client.UpdateTaskList(opts)
		if err != nil {
			t.Fatalf("UpdateTaskList failed: %v", err)
		}
	})

	// 6. Delete the list
	output = captureOutput(t, func() {
		opts := DeleteTaskListOptions{TaskListID: listID}
		err := client.DeleteTaskList(opts)
		if err != nil {
			t.Fatalf("DeleteTaskList failed: %v", err)
		}
	})

	// 7. Final list should not contain the deleted list
	output = captureOutput(t, func() {
		lists, err := client.ListTaskLists()
		if err != nil {
			t.Fatalf("ListTaskLists failed: %v", err)
		}
		printTaskLists(lists)
	})
	if strings.Contains(output, "Groceries") {
		t.Errorf("expected list to not contain 'Groceries', got '%s'", output)
	}
}
