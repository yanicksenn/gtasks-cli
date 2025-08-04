package e2e

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	cliName = "gtasks"
)

func TestMain(m *testing.M) {
	os.Chdir("..")
	os.Exit(m.Run())
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

func TestQuietFlag(t *testing.T) {
	t.Skip("skipping quiet flag e2e tests until authentication is handled in tests")
	cmd := exec.Command("./" + cliName, "tasklists", "list", "--quiet")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run command with --quiet flag: %v\nOutput: %s", err, output)
	}

	if string(output) != "" {
		t.Errorf("expected empty output, got '%s'", output)
	}
}
