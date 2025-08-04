package ui

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"google.golang.org/api/tasks/v1"
	"gopkg.in/yaml.v3"
)

func TestPrinter_PrintTaskLists(t *testing.T) {
	lists := &tasks.TaskLists{
		Items: []*tasks.TaskList{
			{Id: "1", Title: "List 1"},
			{Id: "2", Title: "List 2"},
		},
	}

	t.Run("table", func(t *testing.T) {
		var buf bytes.Buffer
		p := NewPrinter(&buf, "table", false)
		p.PrintTaskLists(lists)
		output := buf.String()
		if !strings.Contains(output, "List 1") {
			t.Errorf("expected output to contain 'List 1', got '%s'", output)
		}
	})

	t.Run("json", func(t *testing.T) {
		var buf bytes.Buffer
		p := NewPrinter(&buf, "json", false)
		p.PrintTaskLists(lists)
		var decoded tasks.TaskLists
		if err := json.NewDecoder(&buf).Decode(&decoded); err != nil {
			t.Fatalf("failed to decode json: %v", err)
		}
		if len(decoded.Items) != 2 {
			t.Errorf("expected 2 items, got %d", len(decoded.Items))
		}
	})

	t.Run("yaml", func(t *testing.T) {
		var buf bytes.Buffer
		p := NewPrinter(&buf, "yaml", false)
		p.PrintTaskLists(lists)
		var decoded tasks.TaskLists
		if err := yaml.NewDecoder(&buf).Decode(&decoded); err != nil {
			t.Fatalf("failed to decode yaml: %v", err)
		}
		if len(decoded.Items) != 2 {
			t.Errorf("expected 2 items, got %d", len(decoded.Items))
		}
	})
}

func TestPrinter_PrintAccounts(t *testing.T) {
	accounts := []string{"one@example.com", "two@example.com"}

	t.Run("table", func(t *testing.T) {
		var buf bytes.Buffer
		p := NewPrinter(&buf, "table", false)
		p.PrintAccounts(accounts, "one@example.com")
		output := buf.String()
		if !strings.Contains(output, "one@example.com (active)") {
			t.Errorf("expected output to contain 'one@example.com (active)', got '%s'", output)
		}
	})

	t.Run("json", func(t *testing.T) {
		var buf bytes.Buffer
		p := NewPrinter(&buf, "json", false)
		p.PrintAccounts(accounts, "")
		var decoded []string
		if err := json.NewDecoder(&buf).Decode(&decoded); err != nil {
			t.Fatalf("failed to decode json: %v", err)
		}
		if len(decoded) != 2 {
			t.Errorf("expected 2 items, got %d", len(decoded))
		}
	})

	t.Run("yaml", func(t *testing.T) {
		var buf bytes.Buffer
		p := NewPrinter(&buf, "yaml", false)
		p.PrintAccounts(accounts, "")
		var decoded []string
		if err := yaml.NewDecoder(&buf).Decode(&decoded); err != nil {
			t.Fatalf("failed to decode yaml: %v", err)
		}
		if len(decoded) != 2 {
			t.Errorf("expected 2 items, got %d", len(decoded))
		}
	})
}
