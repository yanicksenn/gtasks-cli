package gtasks

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestListTaskLists(t *testing.T) {
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

	err = client.ListTaskLists()
	if err != nil {
		t.Fatalf("ListTaskLists failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Test Task List 1") {
		t.Errorf("expected output to contain 'Test Task List 1', got '%s'", output)
	}
	if !strings.Contains(output, "Test Task List 2") {
		t.Errorf("expected output to contain 'Test Task List 2', got '%s'", output)
	}
}

func TestGetTaskList(t *testing.T) {
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

	opts := GetTaskListOptions{TaskListID: "taskList1"}
	err = client.GetTaskList(opts)
	if err != nil {
		t.Fatalf("GetTaskList failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Test Task List 1") {
		t.Errorf("expected output to contain 'Test Task List 1', got '%s'", output)
	}
}

func TestCreateTaskList(t *testing.T) {
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

	opts := CreateTaskListOptions{Title: "New Task List"}
	err = client.CreateTaskList(opts)
	if err != nil {
		t.Fatalf("CreateTaskList failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Successfully created task list: New Task List (newTaskList)") {
		t.Errorf("expected output to contain 'Successfully created task list: New Task List (newTaskList)', got '%s'", output)
	}
}

func TestUpdateTaskList(t *testing.T) {
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

	opts := UpdateTaskListOptions{TaskListID: "taskList1", Title: "Updated Task List"}
	err = client.UpdateTaskList(opts)
	if err != nil {
		t.Fatalf("UpdateTaskList failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Successfully updated task list: Updated Task List (taskList1)") {
		t.Errorf("expected output to contain 'Successfully updated task list: Updated Task List (taskList1)', got '%s'", output)
	}
}

func TestDeleteTaskList(t *testing.T) {
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

	opts := DeleteTaskListOptions{TaskListID: "taskList1"}
	err = client.DeleteTaskList(opts)
	if err != nil {
		t.Fatalf("DeleteTaskList failed: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Successfully deleted task list: taskList1") {
		t.Errorf("expected output to contain 'Successfully deleted task list: taskList1', got '%s'", output)
	}
}
