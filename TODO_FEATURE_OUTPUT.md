# TODO: Feature - Advanced Output Formatting

## 1. Goals

The primary goal of this feature is to enable machine-readable output from the CLI, allowing `gtasks-cli` to be used in scripts and automated workflows.

-   **High-Level Goals:**
    -   Introduce a persistent `--output` (or `-o`) flag to allow users to specify the desired output format.
    -   Support `json` and `yaml` as output formats.
    -   The default behavior (no flag) will remain the current human-readable, table-like format.
    -   Ensure the structured output is well-formed and easy to parse.

-   **Example Usage:**
    ```bash
    # Get a list of tasks in JSON format
    gtasks tasks list --output json

    # Get details for a single task list in YAML format
    gtasks tasklists get <list-id> -o yaml
    ```

## 2. Technical Trade-offs

The implementation for this feature is relatively straightforward, but we should ensure it is done in a clean, maintainable way.

-   **Implementation Strategy:**
    -   **Flag Definition:** A new persistent flag, `--output`, will be added to the `RootCmd` in `cmd/root.go`. This will make it available to all subcommands.
    -   **Centralized Logic:** The logic for handling the different output formats should be centralized within the `internal/ui` package. This avoids duplicating formatting logic across different commands.

-   **Proposed `ui` Package Structure:**
    -   We can create a `Printer` struct in the `ui` package.
    -   This `Printer` will be initialized with the desired output format (e.g., `json`, `yaml`, `table`) and an `io.Writer` (like `os.Stdout`).
    -   The `Printer` will have methods for each data type we need to print, such as:
        -   `PrintTaskList(list *tasks.TaskList)`
        -   `PrintTaskLists(lists *tasks.TaskLists)`
        -   `PrintTask(task *tasks.Task)`
        -   `PrintTasks(tasks *tasks.Tasks)`
    -   Inside these methods, a `switch` statement on the output format will delegate to the appropriate rendering function (e.g., `printTasksAsJSON`, `printTasksAsTable`).

-   **Refactoring `cmd` Package:**
    -   The `CommandHelper` will be updated to initialize this new `ui.Printer`.
    -   The `RunE` function for each command will no longer contain any `fmt.Printf` or `cmd.Println` calls. Instead, it will call the appropriate method on the `helper.Printer`. This further cleans up the command logic and centralizes all output handling.

## 3. Testability

This feature is highly testable and we can achieve excellent test coverage.

-   **Unit Testing the `Printer`:**
    -   We can write dedicated unit tests for the `Printer` in the `ui` package.
    -   For each supported format (`json`, `yaml`, `table`), we will write a test case for each `Print` method.
    -   The testing process will be:
        1.  Create a `Printer` configured for a specific format, using a `bytes.Buffer` as the `io.Writer`.
        2.  Call the method under test with mock `*tasks.Task` or `*tasks.TaskList` data.
        3.  Capture the output from the `bytes.Buffer`.
        4.  **For `json`/`yaml`:** Unmarshal the output string and assert that the resulting data structure is deeply equal to the original mock data. This provides strong verification.
        5.  **For `table`:** Assert that the output string contains expected substrings (e.g., headers, specific task titles).

-   **Integration Testing:**
    -   The existing E2E tests can be extended to include a test case for the `--output json` flag to ensure it is wired up correctly at the command level.
