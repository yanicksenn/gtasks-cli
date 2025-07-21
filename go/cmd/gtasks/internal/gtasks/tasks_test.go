package gtasks

import (
	"strings"
	"testing"
)

func TestTasksLifecycle(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// 1. Create a task list to work with
	var taskListID string
	output := captureOutput(t, func() {
		opts := CreateTaskListOptions{Title: "Shopping List"}
		if err := client.CreateTaskList(opts); err != nil {
			t.Fatalf("CreateTaskList failed: %v", err)
		}
	})
	taskListID = extractID(output)

	// 2. Initial list of tasks should be empty
	output = captureOutput(t, func() {
		opts := ListTasksOptions{TaskListID: taskListID}
		if err := client.ListTasks(opts); err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
	})
	if !strings.Contains(output, "No tasks found.") {
		t.Errorf("expected empty list, got '%s'", output)
	}

	// 3. Create a task
	var taskID string
	output = captureOutput(t, func() {
		opts := CreateTaskOptions{TaskListID: taskListID, Title: "Buy Milk"}
		if err := client.CreateTask(opts); err != nil {
			t.Fatalf("CreateTask failed: %v", err)
		}
	})
	taskID = extractID(output)

	// 4. List should now contain the new task
	output = captureOutput(t, func() {
		opts := ListTasksOptions{TaskListID: taskListID}
		if err := client.ListTasks(opts); err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
	})
	if !strings.Contains(output, "[ ] Buy Milk") {
		t.Errorf("expected list to contain '[ ] Buy Milk', got '%s'", output)
	}

	// 5. Get the task by ID
	output = captureOutput(t, func() {
		opts := GetTaskOptions{TaskListID: taskListID, TaskID: taskID}
		if err := client.GetTask(opts); err != nil {
			t.Fatalf("GetTask failed: %v", err)
		}
	})
	if !strings.Contains(output, "Title:   Buy Milk") {
		t.Errorf("expected get to show 'Title:   Buy Milk', got '%s'", output)
	}

	// 6. Update the task
	output = captureOutput(t, func() {
		opts := UpdateTaskOptions{TaskListID: taskListID, TaskID: taskID, Title: "Buy Almond Milk"}
		if err := client.UpdateTask(opts); err != nil {
			t.Fatalf("UpdateTask failed: %v", err)
		}
	})
	if !strings.Contains(output, "Successfully updated task: Buy Almond Milk") {
		t.Errorf("expected update message, got '%s'", output)
	}

	// 7. Complete the task
	output = captureOutput(t, func() {
		opts := CompleteTaskOptions{TaskListID: taskListID, TaskID: taskID}
		if err := client.CompleteTask(opts); err != nil {
			t.Fatalf("CompleteTask failed: %v", err)
		}
	})
	if !strings.Contains(output, "Successfully completed task: Buy Almond Milk") {
		t.Errorf("expected completion message, got '%s'", output)
	}

	// 8. List should show the task as completed
	output = captureOutput(t, func() {
		opts := ListTasksOptions{TaskListID: taskListID, ShowCompleted: true}
		if err := client.ListTasks(opts); err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
	})
	if !strings.Contains(output, "[x] Buy Almond Milk") {
		t.Errorf("expected list to contain '[x] Buy Almond Milk', got '%s'", output)
	}

	// 9. Delete the task
	output = captureOutput(t, func() {
		opts := DeleteTaskOptions{TaskListID: taskListID, TaskID: taskID}
		if err := client.DeleteTask(opts); err != nil {
			t.Fatalf("DeleteTask failed: %v", err)
		}
	})
	if !strings.Contains(output, "Successfully deleted task") {
		t.Errorf("expected deletion message, got '%s'", output)
	}

	// 10. Final list of tasks should be empty again
	output = captureOutput(t, func() {
		opts := ListTasksOptions{TaskListID: taskListID}
		if err := client.ListTasks(opts); err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
	})
	if !strings.Contains(output, "No tasks found.") {
		t.Errorf("expected empty list, got '%s'", output)
	}
}
