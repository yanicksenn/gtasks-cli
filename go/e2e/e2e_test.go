package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

const (
	cliName = "gtasks"
)

func TestMain(m *testing.M) {
	// Build the CLI before running tests.
	cmd := exec.Command("go", "build", "-o", cliName, "../cmd/gtasks")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(string(output))
	}

	m.Run()
}

func TestHelp(t *testing.T) {
	cmd := exec.Command("./" + cliName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run help command: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Usage:") {
		t.Errorf("expected help text to contain 'Usage:', got '%s'", output)
	}
}

func TestTasklists(t *testing.T) {
	t.Skip("skipping tasklists e2e tests until authentication is handled in tests")
	// TODO: Implement e2e tests for tasklists, which will require handling authentication.
}

func TestTasks(t *testing.T) {
	t.Skip("skipping tasks e2e tests until authentication is handled in tests")
	// TODO: Implement e2e tests for tasks, which will require handling authentication.
}

func TestVersion(t *testing.T) {
	cmd := exec.Command("./" + cliName, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run version command: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "0.1.0") {
		t.Errorf("expected version text to contain '0.1.0', got '%s'", output)
	}
}
