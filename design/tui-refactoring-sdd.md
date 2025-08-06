# TUI Refactoring Software Design Document

**Date:** 2025-08-06

**Goal:** Refactor the existing Bubble Tea TUI to simplify its functionality and align it with the other commands. The new TUI will no longer have a task list selection view. Instead, the user will provide the task list via a command-line flag. The main view will display a list of tasks with their completion status, and users will be able to toggle the completion status of a task.

**Technical Considerations:**

*   **Framework:** We will continue to use the Bubble Tea framework.
*   **Entry Point:** The `gtasks interactive` command will remain the entry point for the TUI.
*   **Task List Selection:** The task list will be provided via a `--tasklist` flag, similar to the `gtasks tasks list` command.
*   **Task Display:** Tasks will be displayed in the format `[ ] Task Title` for incomplete tasks and `[x] Task Title` for completed tasks.
*   **Functionality:**
    *   The TUI will display a list of tasks for a given task list.
    *   Users will be able to navigate the task list using the up and down arrow keys.
    *   Users will be able to toggle the completion status of a task by pressing the spacebar.
    *   The TUI will exit when the user presses `q` or `ctrl+c`.
*   **Removed Features:**
    *   Task list selection view.
    *   Task creation.
    *   Task deletion.
    *   Task detail view.
*   **Dependencies:**
    *   `github.com/charmbracelet/bubbletea`
    *   `github.com/charmbracelet/lipgloss`
    *   `google.golang.org/api/tasks/v1`

**Testing Plan:**

*   **Unit Tests:**
    *   Create unit tests for the new TUI model, including:
        *   Initial model state.
        *   Update logic for navigating the task list.
        *   Update logic for toggling task completion.
        *   View rendering.
*   **E2E Tests:**
    *   Update the existing E2E tests to reflect the new TUI functionality.
    *   Create a new E2E test to verify that the TUI displays the correct tasks for a given task list.
    *   Create a new E2E test to verify that toggling the completion status of a task works correctly.
