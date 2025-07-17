package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func TestTodoCLI(t *testing.T) {
	testdataDir, err := bazel.Runfile("go/cmd/todo/testdata")
	if err != nil {
		t.Fatalf("Could not find testdata directory: %v", err)
	}

	todoBinary, err := bazel.Runfile("go/cmd/todo/todo_/todo")
	if err != nil {
		t.Fatalf("Could not find todo binary: %v", err)
	}

	cases, err := ioutil.ReadDir(testdataDir)
	if err != nil {
		t.Fatalf("Could not read testdata directory: %v", err)
	}

	for _, c := range cases {
		if !c.IsDir() {
			continue
		}

		t.Run(c.Name(), func(t *testing.T) {
			caseDir := filepath.Join(testdataDir, c.Name())
			
			args := []string{}
			if strings.Contains(c.Name(), "aggregated") {
				args = append(args, "--aggregated")
			}
			if strings.Contains(c.Name(), "validation") {
				args = append(args, "--validate")
			}
			args = append(args, caseDir)

			cmd := exec.Command(todoBinary, args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			runErr := cmd.Run()

			expectedOutput, err := ioutil.ReadFile(filepath.Join(caseDir, "expected.out"))
			if err != nil {
				t.Fatalf("Could not read expected.out: %v", err)
			}

			actualOutput := strings.TrimSpace(stdout.String())
			expectedOutputStr := strings.TrimSpace(string(expectedOutput))

			// The tool prints absolute paths, so we need to make them relative for the test
			actualOutput = strings.ReplaceAll(actualOutput, caseDir+"/", "")

			if expectedOutputStr != actualOutput {
				t.Errorf("Output does not match expected.out.\nExpected:\n---\n%s\n---\nActual:\n---\n%s\n---", expectedOutputStr, actualOutput)
			}

			// Handle exit code validation
			if strings.Contains(c.Name(), "validation") {
				shouldFail := strings.Contains(c.Name(), "invalid")
				if shouldFail {
					if runErr == nil {
						t.Fatalf("Expected a non-zero exit code for validation test, but got nil")
					}
				} else {
					if runErr != nil {
						t.Fatalf("Expected a zero exit code for validation test, but got: %v", runErr)
					}
				}
			} else {
				if runErr != nil {
					t.Fatalf("Command failed with error: %v\nStderr: %s", runErr, stderr.String())
				}
			}
		})
	}
}