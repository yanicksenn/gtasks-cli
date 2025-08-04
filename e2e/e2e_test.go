package e2e

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/yanicksenn/gtasks/cmd"
	"github.com/yanicksenn/gtasks/internal/version"
)

func execute(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetErr(buf)
	cmd.RootCmd.SetArgs(args)
	err := cmd.RootCmd.Execute()
	return buf.String(), err
}

func login() error {
	_, err := execute("accounts", "login")
	return err
}

func TestMain(m *testing.M) {
	// Setup: login
	if err := login(); err != nil {
		fmt.Fprintf(os.Stderr, "e2e setup: failed to login: %v\n", err)
		os.Exit(1)
	}

	exitVal := m.Run()

	// Teardown: logout
	_, err := execute("accounts", "logout")
	if err != nil {
		fmt.Fprintf(os.Stderr, "e2e teardown: failed to logout: %v\n", err)
	}
	os.Exit(exitVal)
}

func TestHelp(t *testing.T) {
	output, err := execute()
	if err != nil {
		t.Fatalf("failed to run help command: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("expected help text to contain 'Usage:', got '%s'", output)
	}
	if !strings.Contains(output, cmd.RootCmd.Long) {
		t.Errorf("expected help text to contain the long description, got '%s'", output)
	}
	if !strings.Contains(output, "Available Commands:") {
		t.Errorf("expected help text to contain 'Available Commands:', got '%s'", output)
	}
}

func TestTasklists(t *testing.T) {
	// Create a new task list
	listTitle := "E2E Test List"
	output, err := execute("tasklists", "create", "--title", listTitle)
	if err != nil {
		t.Fatalf("failed to create task list: %v\nOutput: %s", err, output)
	}

	// Extract the ID from the output
	re := regexp.MustCompile(`Successfully created task list: .* \((.*)\)`) 
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		t.Fatalf("could not find task list ID in output: %s", output)
	}
	listID := matches[1]

	// Ensure the task list is deleted even if the test fails
	t.Cleanup(func() {
		_, err := execute("tasklists", "delete", listID)
		if err != nil {
			t.Logf("failed to delete task list %s during cleanup: %v", listID, err)
		}
	})

	// List task lists and verify the new one is there
	output, err = execute("tasklists", "list")
	if err != nil {
		t.Fatalf("failed to list task lists: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, listTitle) {
		t.Errorf("expected task list '%s' to be in the list", listTitle)
	}

	// Get the task list and verify its details
	output, err = execute("tasklists", "get", listID)
	if err != nil {
		t.Fatalf("failed to get task list: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, listTitle) {
		t.Errorf("expected task list details to contain '%s'", listTitle)
	}
}

func TestTasks(t *testing.T) {
	// Create a new task list for this test
	listTitle := "E2E Tasks Test List"
	output, err := execute("tasklists", "create", "--title", listTitle)
	if err != nil {
		t.Fatalf("failed to create task list for tasks test: %v\nOutput: %s", err, output)
	}
	re := regexp.MustCompile(`Successfully created task list: .* \((.*)\)`) 
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		t.Fatalf("could not find task list ID in output: %s", output)
	}
	listID := matches[1]
	t.Cleanup(func() {
		_, err := execute("tasklists", "delete", listID)
		if err != nil {
			t.Logf("failed to delete task list %s during cleanup: %v", listID, err)
		}
	})

	// Create a new task
	taskTitle := "E2E Test Task"
	output, err = execute("tasks", "create", "--tasklist", listID, "--title", taskTitle)
	if err != nil {
		t.Fatalf("failed to create task: %v\nOutput: %s", err, output)
	}
	re = regexp.MustCompile(`Successfully created task: .* \((.*)\)`) 
	matches = re.FindStringSubmatch(output)
	if len(matches) < 2 {
		t.Fatalf("could not find task ID in output: %s", output)
	}
	taskID := matches[1]

	// List tasks and verify the new one is there
	output, err = execute("tasks", "list", "--tasklist", listID)
	if err != nil {
		t.Fatalf("failed to list tasks: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, taskTitle) {
		t.Errorf("expected task '%s' to be in the list", taskTitle)
	}

	// Get the task and verify its details
	output, err = execute("tasks", "get", "--tasklist", listID, taskID)
	if err != nil {
		t.Fatalf("failed to get task: %v\nOutput: %s", err, output)
	}
	if !strings.Contains(output, taskTitle) {
		t.Errorf("expected task details to contain '%s'", taskTitle)
	}
}

func TestVersion(t *testing.T) {
	output, err := execute("--version")
	if err != nil {
		t.Fatalf("failed to run version command: %v\nOutput: %s", err, output)
	}

	if strings.TrimSpace(output) != version.Get() {
		t.Errorf("expected version text to be '%s', got '%s'", version.Get(), output)
	}
}

func TestQuietFlag(t *testing.T) {
	output, err := execute("tasklists", "list", "--quiet")
	if err != nil {
		t.Fatalf("failed to run command with --quiet flag: %v\nOutput: %s", err, output)
	}

	if strings.TrimSpace(output) != "" {
		t.Errorf("expected empty output, got '%s'", output)
	}
}
