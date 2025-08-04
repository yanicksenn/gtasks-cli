# Design Document: Sorting Task Lists and Tasks

## 1. Goals

The primary goal of this feature is to provide users with the ability to sort their task lists and tasks based on various criteria, allowing for better organization and prioritization.

## 2. Sorting Options

### 2.1. Task Lists

The `tasklists list` command will support the following sorting options via a `--sort-by` flag:

-   **`alphabetical` (default):** Sorts task lists by their title in ascending alphabetical order.
-   **`last-modified`:** Sorts task lists by their last modified timestamp in descending order (newest first).
-   **`uncompleted-tasks`:** Sorts task lists by the number of uncompleted tasks they contain, in descending order (most uncompleted tasks first).

### 2.2. Tasks

The `tasks list` command will support the following sorting options via a `--sort-by` flag:

-   **`alphabetical` (default):** Sorts tasks by their title in ascending alphabetical order.
-   **`last-modified`:** Sorts tasks by their last modified timestamp in descending order (newest first).
-   **`due-date`:** Sorts tasks by their due date in ascending order (earliest first). Tasks without a due date will be placed at the end. Tasks with the same due date will be sorted alphabetically by title.

## 3. Implementation Details

### 3.1. CLI

-   A `--sort-by` flag will be added to the `tasklists list` and `tasks list` commands.
-   The flag will accept one of the predefined sorting options for the respective command.
-   The default sorting option for both commands will be `alphabetical`.

### 3.2. `gtasks` Client

-   The `ListTaskListsOptions` and `ListTasksOptions` structs will be updated to include a `SortBy` field.
-   The `ListTaskLists` and `ListTasks` methods will be updated to perform the sorting based on the `SortBy` field.
-   For the `uncompleted-tasks` sorting option, the `ListTaskLists` method will need to fetch the tasks for each task list to get the uncompleted task count. This will be a performance consideration, and the implementation should be mindful of potential rate limiting.

### 3.3. TUI

-   A keybinding (e.g., `s`) will be added to the TUI to cycle through the available sorting options for the focused pane.
-   The current sorting option will be displayed in the status bar.
-   When the sorting option is changed, the list will be re-fetched and re-rendered with the new sorting.

## 4. Testing

-   Unit tests will be added for the sorting logic in the `gtasks` client.
-   E2E tests will be added for the `--sort-by` flag in the CLI commands.

## 5. Documentation

-   The `README.md` will be updated to document the new `--sort-by` flag for the `tasklists list` and `tasks list` commands.
-   The `README.md` will be updated to document the new keybinding for sorting in the interactive mode.
