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
  - [ ] **3.1:** Define a `Client` interface in `internal/gtasks/client.go` that represents all high-level operations.
  - [ ] **3.2:** Rename the existing `Client` struct to `onlineClient` and make it implement the `Client` interface.
  - [ ] **3.3:** Refactor all methods on the `onlineClient` to **return** data (e.g., `*tasks.TaskList`) instead of printing it.
  - [ ] **3.4:** Update the `Run` functions in the `cmd/` package to handle all printing based on the data returned from the client.

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