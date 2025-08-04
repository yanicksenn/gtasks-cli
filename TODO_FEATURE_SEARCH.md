# TODO: Feature - Advanced Search and Filtering

## 1. Goals

The goal of this feature is to provide users with the ability to find specific tasks by searching and filtering, rather than manually scanning through long lists.

-   **High-Level Goals:**
    -   Allow users to filter tasks based on their properties, such as title, notes, and due date.
    -   The filtering should be performed on the client side, as the Google Tasks API does not support server-side searching on these fields.
    -   Integrate the filtering functionality seamlessly into the existing `tasks list` command.

-   **Proposed Filtering Flags:**
    -   `--title-contains <string>`: Filter tasks where the title contains the given substring (case-insensitive).
    -   `--notes-contains <string>`: Filter tasks where the notes contain the given substring (case-insensitive).
    -   `--due-before <date>`: Filter tasks with a due date before the specified date (e.g., "2025-12-31").
    -   `--due-after <date>`: Filter tasks with a due date after the specified date.

-   **Example Usage:**
    ```bash
    # Find all tasks with "meeting" in the title
    gtasks tasks list --title-contains "meeting"

    # Find all tasks due before the end of the year
    gtasks tasks list --due-before 2025-12-31
    ```

## 2. Technical Trade-offs

-   **API Limitations:**
    -   The Google Tasks API is limited. It does not offer a search endpoint or server-side filtering on most task fields.
    -   This means our implementation **must** be client-side. The workflow for a filtered list will be:
        1.  Fetch *all* tasks from the specified task list using the existing `client.ListTasks` method.
        2.  Apply the user-provided filters to this list of tasks in memory.
        3.  Display only the tasks that match the criteria.
    -   **Performance Consideration:** For users with extremely large task lists (thousands of tasks), this could introduce a small amount of latency, but for the vast majority of users, it will be unnoticeable.

-   **Command Structure:**
    -   **Option A: Add flags to `tasks list` (Recommended)**
        -   **Pros:** This is the most intuitive approach. Filtering is a natural extension of listing. It keeps the command structure simple and avoids adding a new top-level command.
        -   **Cons:** If we add too many filter flags, the `tasks list` command could become bloated.
    -   **Option B: Create a new `tasks find` command**
        -   **Pros:** Creates a clear separation between simply listing all items and performing a filtered search.
        -   **Cons:** Adds another command to maintain and might be overkill if the filtering options are not overly complex.

-   **Implementation Plan (using Option A):**
    1.  Add the new filtering flags (`--title-contains`, etc.) to the `listTasksCmd` in `cmd/tasks.go`.
    2.  In the `RunE` function for the command, parse these flags.
    3.  Create a new `gtasks.FilterOptions` struct to hold the parsed filter values.
    4.  The core filtering logic will be implemented in a new function, likely `gtasks.FilterTasks(tasks []*tasks.Task, opts FilterOptions) []*tasks.Task`. This function will be pure and easily testable.
    5.  The `RunE` function will call `h.Client.ListTasks`, and then pass the result through the `gtasks.FilterTasks` function before printing the output.

## 3. Testability

This feature is highly testable, especially the core filtering logic.

-   **Unit Testing the Filter Logic:**
    -   The `gtasks.FilterTasks` function will be the primary target for unit tests.
    -   We will create a suite of tests that:
        1.  Defines a static, mock list of `*tasks.Task` objects with varying titles, notes, and due dates.
        2.  Calls `FilterTasks` with different `FilterOptions` (e.g., filtering by title, by date, by a combination of both).
        3.  Asserts that the returned slice of tasks is correct (i.e., it contains only the expected tasks and the count is correct).
        4.  We will test edge cases, such as case-insensitivity, empty filter criteria, and no matching tasks.

-   **Integration Testing:**
    -   We can add a new test case to `e2e/e2e_test.go` to verify that the flags are correctly wired up.
    -   The test would:
        1.  Create a task list.
        2.  Create a few tasks with specific, known titles.
        3.  Run `gtasks tasks list --title-contains "unique-keyword"`.
        4.  Assert that the output contains the task with the unique keyword and not the others.
