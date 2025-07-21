package gtasks

import (
	"bytes"
	"io"
	"os"
	"regexp"
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
		if err := client.ListTaskLists(); err != nil {
			t.Fatalf("ListTaskLists failed: %v", err)
		}
	})
	if !strings.Contains(output, "Default List") {
		t.Errorf("expected default list, got '%s'", output)
	}

	// 2. Create a new task list
	var listID string
	output = captureOutput(t, func() {
		opts := CreateTaskListOptions{Title: "Groceries"}
		if err := client.CreateTaskList(opts); err != nil {
			t.Fatalf("CreateTaskList failed: %v", err)
		}
	})
	if !strings.Contains(output, "Successfully created task list: Groceries") {
		t.Errorf("expected creation message, got '%s'", output)
	}
	listID = extractID(output)

	// 3. List should now contain the new list
	output = captureOutput(t, func() {
		if err := client.ListTaskLists(); err != nil {
			t.Fatalf("ListTaskLists failed: %v", err)
		}
	})
	if !strings.Contains(output, "Groceries") {
		t.Errorf("expected list to contain 'Groceries', got '%s'", output)
	}

	// 4. Get the list by ID
	output = captureOutput(t, func() {
		opts := GetTaskListOptions{TaskListID: listID}
		if err := client.GetTaskList(opts); err != nil {
			t.Fatalf("GetTaskList failed: %v", err)
		}
	})
	if !strings.Contains(output, "Title: Groceries") {
		t.Errorf("expected get to show 'Title: Groceries', got '%s'", output)
	}

	// 5. Update the list
	output = captureOutput(t, func() {
		opts := UpdateTaskListOptions{TaskListID: listID, Title: "Updated Groceries"}
		if err := client.UpdateTaskList(opts); err != nil {
			t.Fatalf("UpdateTaskList failed: %v", err)
		}
	})
	if !strings.Contains(output, "Successfully updated task list: Updated Groceries") {
		t.Errorf("expected update message, got '%s'", output)
	}

	// 6. Delete the list
	output = captureOutput(t, func() {
		opts := DeleteTaskListOptions{TaskListID: listID}
		if err := client.DeleteTaskList(opts); err != nil {
			t.Fatalf("DeleteTaskList failed: %v", err)
		}
	})
	if !strings.Contains(output, "Successfully deleted task list") {
		t.Errorf("expected deletion message, got '%s'", output)
	}

	// 7. Final list should not contain the deleted list
	output = captureOutput(t, func() {
		if err := client.ListTaskLists(); err != nil {
			t.Fatalf("ListTaskLists failed: %v", err)
		}
	})
	if strings.Contains(output, "Groceries") {
		t.Errorf("expected list to not contain 'Groceries', got '%s'", output)
	}
}

// captureOutput is a helper to capture stdout from a function.
func captureOutput(t *testing.T, f func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// extractID is a helper to extract an ID from the CLI's output.
func extractID(output string) string {
	re := regexp.MustCompile(`\((.*?)\)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}