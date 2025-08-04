package gtasks

import (
	"testing"

	taskspb "google.golang.org/api/tasks/v1"
)

func TestFilterTasks(t *testing.T) {
	tasks := []*taskspb.Task{
		{
			Title: "Buy milk",
			Notes: "Get the good stuff",
			Due:   "2025-12-20T15:00:00.000Z",
		},
		{
			Title: "Finish report",
			Notes: "For the Q4 meeting",
			Due:   "2025-12-22T18:00:00.000Z",
		},
		{
			Title: "Call mom",
			Notes: "Wish her a happy birthday",
			Due:   "2025-12-25T10:00:00.000Z",
		},
	}

	testCases := []struct {
		name          string
		opts          FilterOptions
		expectedCount int
		expectedTitle string
	}{
		{
			name:          "Filter by title",
			opts:          FilterOptions{TitleContains: "milk"},
			expectedCount: 1,
			expectedTitle: "Buy milk",
		},
		{
			name:          "Filter by notes",
			opts:          FilterOptions{NotesContains: "meeting"},
			expectedCount: 1,
			expectedTitle: "Finish report",
		},
		{
			name:          "Filter by due before",
			opts:          FilterOptions{DueBefore: "2025-12-21"},
			expectedCount: 1,
			expectedTitle: "Buy milk",
		},
		{
			name:          "Filter by due after",
			opts:          FilterOptions{DueAfter: "2025-12-24"},
			expectedCount: 1,
			expectedTitle: "Call mom",
		},
		{
			name:          "No filters",
			opts:          FilterOptions{},
			expectedCount: 3,
		},
		{
			name:          "No matching tasks",
			opts:          FilterOptions{TitleContains: "nonexistent"},
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filtered, err := FilterTasks(tasks, tc.opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(filtered) != tc.expectedCount {
				t.Errorf("expected %d tasks, but got %d", tc.expectedCount, len(filtered))
			}

			if tc.expectedCount == 1 && filtered[0].Title != tc.expectedTitle {
				t.Errorf("expected task with title '%s', but got '%s'", tc.expectedTitle, filtered[0].Title)
			}
		})
	}
}