package gtasks

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestListTasks(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := ListTasksOptions{TaskListID: "taskList1"}
	err = client.ListTasks(opts)
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Test Task 1") {
		t.Errorf("expected output to contain 'Test Task 1', got '%s'", output)
	}
	if !strings.Contains(output, "Test Task 2") {
		t.Errorf("expected output to contain 'Test Task 2', got '%s'", output)
	}
}

func TestGetTask(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := GetTaskOptions{TaskListID: "taskList1", TaskID: "task1"}
	err = client.GetTask(opts)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Test Task 1") {
		t.Errorf("expected output to contain 'Test Task 1', got '%s'", output)
	}
}

func TestCreateTask(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := CreateTaskOptions{TaskListID: "taskList1", Title: "New Task"}
	err = client.CreateTask(opts)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Successfully created task: New Task (newTask)") {
		t.Errorf("expected output to contain 'Successfully created task: New Task (newTask)', got '%s'", output)
	}
}

func TestUpdateTask(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := UpdateTaskOptions{TaskListID: "taskList1", TaskID: "task1", Title: "Updated Task"}
	err = client.UpdateTask(opts)
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Successfully updated task: Updated Task (task1)") {
		t.Errorf("expected output to contain 'Successfully updated task: Updated Task (task1)', got '%s'", output)
	}
}

func TestCompleteTask(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := CompleteTaskOptions{TaskListID: "taskList1", TaskID: "task1"}
	err = client.CompleteTask(opts)
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Successfully completed task: Updated Task (task1)") {
		t.Errorf("expected output to contain 'Successfully completed task: Updated Task (task1)', got '%s'", output)
	}
}

func TestDeleteTask(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	client, err := newTestClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := DeleteTaskOptions{TaskListID: "taskList1", TaskID: "task1"}
	err = client.DeleteTask(opts)
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Successfully deleted task: task1") {
		t.Errorf("expected output to contain 'Successfully deleted task: task1', got '%s'", output)
	}
}