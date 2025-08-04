# TODO: Refactoring the `gtasks` Package

## Problem

The `internal/gtasks` package currently has mixed responsibilities. It is responsible for:

1.  **API Client Logic:** Defining the `Client` interface and implementing the `onlineClient` for communicating with the Google Tasks API.
2.  **Offline Data Storage:** Containing the `offline_store.go` and `offline_client.go`, which manage local data persistence for the `--offline` mode.
3.  **User-Facing Output:** Including the `print.go` file, which formats and prints data structures like `*tasks.TaskList` directly to the console.

## Why It's Bad

This design violates several core software engineering principles:

-   **Single Responsibility Principle (SRP):** A package should have one, and only one, reason to change. The `gtasks` package currently has three reasons to change: modifications to the Google Tasks API, changes to the offline storage mechanism, or adjustments to the command-line output format. This makes the package brittle and difficult to manage.
-   **Poor Testability:** Unit testing the API client logic in isolation is complicated. Tests for the client may inadvertently depend on the printing or storage logic, making them harder to write and less focused.
-   **Low Cohesion:** The components within the package are not strongly related. For example, `print.go` is fundamentally a UI concern, while `client.go` is about data fetching and remote communication. Mixing these leads to a confusing package structure.
-   **Difficult Maintenance:** When responsibilities are tangled, it becomes challenging for developers to locate the correct place to make changes or fix bugs. This increases the cognitive load and slows down development.

## Proposed Solution

To address these issues, we will refactor the codebase to separate these distinct concerns into their own dedicated packages.

1.  **Create a new `internal/ui` package:**
    *   Move the printing logic from `internal/gtasks/print.go` and its corresponding test, `internal/gtasks/print_test.go`, into this new `ui` package.
    *   The `ui` package will be solely responsible for all user-facing output. It will accept data structures (e.g., `*tasks.Task`, `*tasks.TaskList`) as input and handle their presentation on the command line.

2.  **Create a new `internal/store` package:**
    *   Move the offline storage logic from `internal/gtasks/offline_store.go` and `internal/gtasks/offline_client.go` into this new package.
    *   This move will be part of a larger effort to unify the duplicated storage logic, as detailed in `TODO_REFACTOR_STORE.md`.

3.  **Refine the `gtasks` package:**
    *   After moving the printing and storage code, the `gtasks` package will be left with its core responsibility: defining the `Client` interface and implementing the `onlineClient` that communicates with the live Google Tasks API.
    *   This will result in a lean, focused package that is easy to understand, test, and maintain.
