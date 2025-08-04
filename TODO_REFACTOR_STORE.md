# TODO: Unifying In-Memory Storage Logic

## Problem

The codebase currently contains two nearly identical in-memory storage implementations, which is a significant source of technical debt:

1.  **`internal/gtasks/offline_store.go`**: This implementation provides a local, file-backed persistence layer that allows the CLI to function in `--offline` mode.
2.  **`internal/gtasks/mock_store_test.go`**: This implementation serves as an in-memory mock database for the test suite, particularly for the tests in `tasklists_test.go` and `tasks_test.go`.

A side-by-side comparison of these files reveals that their logic for creating, reading, updating, and deleting tasks and task lists is functionally the same.

## Why It's Bad

This duplication has several negative consequences:

-   **Violation of DRY Principle (Don't Repeat Yourself):** It's a classic code smell. Duplicating logic bloats the codebase and increases the surface area for bugs.
-   **High Maintenance Overhead:** Every time a change or bug fix is required in the storage logic, it must be implemented in two separate places. This is inefficient and error-prone, as it's easy for the two implementations to diverge, leading to inconsistent behavior between the tests and the actual offline application.
-   **Increased Complexity:** The presence of redundant files makes the project harder to navigate and understand for both current and future developers. It creates unnecessary clutter and confusion about which implementation to use where.

## Proposed Solution

We will eliminate this redundancy by creating a single, reusable in-memory store in a new `internal/store` package.

1.  **Create a new `internal/store` package.**
    *   This package will become the canonical location for all data storage-related logic.

2.  **Implement a single `InMemoryStore` in `internal/store/in_memory_store.go`.**
    *   This new implementation will merge the logic from both `offline_store.go` and `mock_store_test.go`.
    *   It will be designed as a flexible, thread-safe, in-memory data store for tasks and task lists.
    *   It will include an optional persistence mechanism (writing to a JSON file) that can be enabled for offline use or disabled for testing, controlled via its constructor.

3.  **Refactor the `offlineClient` to use the new store.**
    *   The `offlineClient` (from `internal/gtasks/offline_client.go`) will be updated to delegate all its operations to an instance of the new `InMemoryStore`. It will no longer contain any storage logic itself.

4.  **Update all relevant tests.**
    *   The unit and integration tests in the `gtasks` package (e.g., `tasklists_test.go`, `tasks_test.go`) will be refactored to use the new `InMemoryStore` directly for mocking, replacing the need for `mock_store_test.go` and the mock server. This will simplify the test setup and make tests run faster.

5.  **Delete the redundant files.**
    *   Once the refactoring is complete and all tests are passing, the following files will be safely removed:
        *   `internal/gtasks/offline_store.go`
        *   `internal/gtasks/offline_store_test.go`
        *   `internal/gtasks/mock_store_test.go`
        *   `internal/gtasks/mock_server_test.go`
