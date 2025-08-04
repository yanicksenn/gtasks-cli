# TODO: Feature - Advanced Output Formatting

## 1. Goals

The primary goal of this feature is to enable machine-readable output from the CLI, allowing `gtasks-cli` to be used in scripts and automated workflows. This involves creating a flexible output system and consolidating the existing `print` and `get` commands into a single, consistent API.

-   **High-Level Goals:**
    -   Introduce a persistent `--output` (or `-o`) flag to allow users to specify the desired output format.
    -   Support `json` and `yaml` as output formats. The default will remain the human-readable `table` format.
    -   **Deprecate the `print` subcommands** (`tasks print` and `tasklists print`) in favor of using the `get` command with the `--output` flag and piping to tools like `jq`. This will create a cleaner and more powerful API.
    -   Ensure the output system can format all relevant data types, including task lists, tasks, and accounts.

-   **Example Usage (New API):**
    ```bash
    # Get a list of tasks in JSON format
    gtasks tasks list --output json

    # Get details for a single task list in YAML format
    gtasks tasklists get <list-id> -o yaml

    # Get just the title of a task (replaces `tasks print`)
    gtasks tasks get <task-id> --output json | jq -r '.title'

    # List all authenticated accounts as a JSON array
    gtasks accounts list --output json
    ```

## 2. Merging `print` and Consolidating the API

The existence of both a `get` command and a `print` command creates a confusing and redundant API. The `print` command is essentially a less powerful version of what `get` combined with standard shell tools can achieve.

-   **The Problem:**
    -   `get`: Fetches and displays a whole object (e.g., a task).
    -   `print`: Fetches a whole object and displays just one property of it.
    -   Having both is unnecessary. A user wanting a single property from a structured format is better served by a dedicated tool like `jq` (for JSON) or `yq` (for YAML). This is a standard and powerful pattern in modern CLIs.

-   **The Solution: Deprecate `print`**
    1.  We will remove the `printTaskCmd` and `printTaskListCmd` from `cmd/tasks.go` and `cmd/tasklists.go` respectively.
    2.  The `get` command will become the sole method for retrieving a single item.
    3.  The `--output` flag will control how that item is rendered (table, JSON, YAML).
    4.  The documentation (e.g., command help text) will be updated to guide users to the new `get ... --output json | jq` pattern.

This change simplifies our API, reduces code maintenance, and empowers users by integrating with the broader CLI ecosystem.

## 3. Technical Implementation

-   **Flag Definition:** A new persistent flag, `--output`, will be added to the `RootCmd` in `cmd/root.go`.

-   **Centralized `ui.Printer`:**
    -   The `internal/ui` package will contain a `Printer` struct responsible for all output.
    -   It will be initialized with the desired output format (`json`, `yaml`, `table`) and an `io.Writer`.
    -   The `Printer` will have methods for each data type:
        -   `PrintTaskList(*tasks.TaskList)`
        -   `PrintTaskLists(*tasks.TaskLists)`
        -   `PrintTask(*tasks.Task)`
        -   `PrintTasks(*tasks.Tasks)`
        -   **`PrintAccounts([]string)`**: A new method to handle printing the list of account emails. In JSON/YAML, this will be a simple array of strings.

-   **Refactoring `cmd` Package:**
    -   The `CommandHelper` will be updated to initialize the `ui.Printer`.
    -   All `RunE` functions will delegate all output to the `helper.Printer`.
    -   The `print` commands will be removed entirely.

## 4. Testability

This feature remains highly testable.

-   **Unit Testing the `Printer`:**
    -   We will add unit tests for the new `PrintAccounts` method for all supported formats.
    -   Existing tests for the `Printer` will be updated to reflect the consolidation of `get` and `print`.
    -   We will test that calling a `Print` method with mock data and a `bytes.Buffer` produces a correctly formatted string for JSON, YAML, and the default table view.

-   **E2E Testing:**
    -   The E2E tests will be updated to remove any tests for the now-defunct `print` command.
    -   New E2E tests will be added to verify the `--output` flag works correctly for `tasks list`, `tasks get`, `tasklists list`, `tasklists get`, and `accounts list`.