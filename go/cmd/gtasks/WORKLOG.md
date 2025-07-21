# Worklog

This document provides a high-level, chronological summary of the major features and changes implemented in the `gtasks` CLI project.

## Initial Implementation

-   **Foundation & Scaffolding:** Set up the basic project structure, Go module, and initial dependencies (Cobra, Google API clients).
-   **Core Authentication:** Implemented the full OAuth2 flow for Google Sign-In, including secure local credential caching.
-   **TaskList Management:** Implemented full CRUD (Create, Read, Update, Delete) operations for `tasklists`.
-   **Task Management:** Implemented full CRUD operations for `tasks`, including a `complete` action.
-   **Account Management:** Added support for multiple accounts, including `login`, `logout`, `list`, and `switch` commands.
-   **Default TaskList:** Updated all task-related commands to use the `@default` task list when no specific list is provided, improving usability.
-   **End-to-End Testing:** Established a basic E2E testing framework to verify the compiled binary.

## Testing Overhaul

-   **Stateful Mock API:** Replaced simple, stateless test mocks with a high-fidelity, stateful, in-memory mock of the Google Tasks API.
-   **Lifecycle Tests:** Implemented comprehensive, mock-based lifecycle tests for both `tasklists` and `tasks`, verifying the entire create-read-update-delete flow without requiring network access.

## Offline Mode Feature

-   **Persistent Offline Store:** Created a file-based, persistent store (`offline.json`) to save all tasks and task lists created in offline mode.
-   **Service Abstraction (Facade):** Refactored the core logic to use a `Client` interface, cleanly separating the command layer from the business logic. This facade allows the CLI to seamlessly switch between online and offline data sources.
-   **`--offline` Flag:** Implemented a global `--offline` flag to enable offline operations.
-   **Offline Client:** Created an `offlineClient` that implements the `Client` interface by interacting with the local JSON file instead of the live Google API.
-   **Removed `sync` Command:** The initial, destructive `sync` command was removed from the codebase and documentation to prevent accidental data loss.
