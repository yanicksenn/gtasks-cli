package ui

import (
	"strings"
	"testing"

	"google.golang.org/api/tasks/v1"
)

func TestPrintTaskListProperty(t *testing.T) {
	list := &tasks.TaskList{
		Id:    "test-id",
		Title: "Test List",
	}

	testCases := []struct {
		property string
		expected string
	}{
		{"id", "test-id"},
		{"title", "Test List"},
		{"ID", "test-id"},
		{"Title", "Test List"},
	}

	for _, tc := range testCases {
		t.Run(tc.property, func(t *testing.T) {
			output := CaptureOutput(t, func() {
				err := PrintTaskListProperty(list, tc.property, false)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			})
			if !strings.Contains(output, tc.expected) {
				t.Errorf("expected output to contain '%s', got '%s'", tc.expected, output)
			}
		})
	}
}

func TestPrintTaskProperty(t *testing.T) {
	task := &tasks.Task{
		Id:     "test-id",
		Title:  "Test Task",
		Notes:  "Test Notes",
		Due:    "2025-12-31T23:59:59.000Z",
		Status: "needsAction",
	}

	testCases := []struct {
		property string
		expected string
	}{
		{"id", "test-id"},
		{"title", "Test Task"},
		{"notes", "Test Notes"},
		{"due", "2025-12-31T23:59:59.000Z"},
		{"status", "needsAction"},
	}

	for _, tc := range testCases {
		t.Run(tc.property, func(t *testing.T) {
			output := CaptureOutput(t, func() {
				err := PrintTaskProperty(task, tc.property, false)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			})
			if !strings.Contains(output, tc.expected) {
				t.Errorf("expected output to contain '%s', got '%s'", tc.expected, output)
			}
		})
	}
}

func TestPrintUnknownProperty(t *testing.T) {
	list := &tasks.TaskList{}
	err := PrintTaskListProperty(list, "unknown", false)
	if err == nil {
		t.Error("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown property: unknown. Available properties: id, title, selfLink") {
		t.Errorf("expected error message to contain 'unknown property: unknown. Available properties: id, title, selfLink', got '%s'", err.Error())
	}

	task := &tasks.Task{}
	err = PrintTaskProperty(task, "unknown", false)
	if err == nil {
		t.Error("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown property: unknown. Available properties: id, title, notes, due, status, selfLink") {
		t.Errorf("expected error message to contain 'unknown property: unknown. Available properties: id, title, notes, due, status, selfLink', got '%s'", err.Error())
	}
}


