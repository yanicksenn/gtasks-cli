# Offline Mode - Implementation Milestones

This document outlines the development plan for the offline mode feature. We will track our progress by checking off each item as it's completed and tested.

- [x] **Milestone 1: Persistent Offline Store**
  - [x] Create a new `internal/gtasks/offline_store.go` file.
  - [x] Implement an `offlineStore` struct that reads from and writes to a persistent JSON file.
  - [x] Add methods to the `offlineStore` for all CRUD operations.
  - [x] Add unit tests for the `offlineStore`.

- [x] **Milestone 2: Offline Mode Switching**
  - [x] Add a global, persistent `--offline` flag to the root command.
  - [x] Create the initial structure for the `NewClient` factory.

- [ ] **Milestone 3: Refactor for Abstraction (The Service Facade)**
  - [ ] **Sub-milestone 3.1: Create the Service Facade**
    - [ ] In `internal/gtasks/client.go`, define the new `Client` interface with all the high-level operation methods (e.g., `ListTaskLists`, `CreateTask`).
    - [ ] Rename the existing `Client` struct to `onlineClient`.
    - [ ] Update the `newOnlineClient` function to return this `*onlineClient` struct.
  - [ ] **Sub-milestone 3.2: Refactor Logic to Return Data**
    - [ ] Go through every method on the `onlineClient` in `tasklists.go` and `tasks.go`.
    - [ ] Change each method's signature to **return** the relevant data structure (e.g., `*tasks.TaskLists`, `*tasks.Task`) and an `error`.
    - [ ] **Crucially, remove all `fmt.Printf` calls from this layer.** The business logic should not be responsible for presentation.
  - [ ] **Sub-milestone 3.3: Update Commands to Handle Presentation**
    - [ ] Go through every `Run` function in `cmd/tasklists.go` and `cmd/tasks.go`.
    - [ ] Update them to call the refactored client methods.
    - [ ] Add logic to handle the returned data and errors.
    - [ ] **Re-implement all the `fmt.Printf` calls here**, in the command layer, to display the results to the user.
  - [ ] **Sub-milestone 3.4: Verify the Refactoring**
    - [ ] Run the entire test suite (`go test ./...`) to confirm that our mock-based tests still pass. This will prove that the refactoring was successful and did not break any existing functionality.

- [ ] **Milestone 4: Implement the Offline Service**
  - [ ] **4.1:** Create a new `offlineClient` struct that also implements the `Client` interface.
  - [ ] **4.2:** Implement the methods of `offlineClient` to call the `offlineStore`.
  - [ ] **4.3:** Update the `NewClient` factory to return either an `onlineClient` or `offlineClient` based on the `--offline` flag.

- [ ] **Milestone 5: The `sync` Command**
  - [ ] Create a new `gtasks sync` command.
  - [ ] Implement the logic to read from the offline store and push changes to the online service.
  - [ ] Add a confirmation prompt before syncing.

- [ ] **Milestone 6: Documentation**
  - [ ] Update the main `README.md` to fully document the `--offline` flag and the new `sync` command.
  - [ ] Update the `DESIGN.md` to reflect the new offline architecture.
