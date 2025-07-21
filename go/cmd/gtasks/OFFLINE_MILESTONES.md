# Offline Mode - Implementation Milestones

This document outlines the development plan for the offline mode feature. We will track our progress by checking off each item as it's completed and tested.

- [ ] **Milestone 1: Persistent Offline Store**
  - [ ] Create a new `internal/gtasks/offline_store.go` file.
  - [ ] Implement an `offlineStore` struct that reads from and writes to a persistent JSON file (e.g., `~/.config/gtasks/offline.json`).
  - [ ] Add methods to the `offlineStore` for all CRUD operations on both `tasklists` and `tasks` (e.g., `createTaskList`, `getTask`).
  - [ ] Add unit tests for the `offlineStore` to ensure data is persisted correctly.

- [ ] **Milestone 2: Offline Mode Switching**
  - [ ] Add a global, persistent `--offline` flag to the root command in `cmd/root.go`.
  - [ ] Refactor the `gtasks.NewClient` function into a "factory" that checks for the `--offline` flag.
  - [ ] The factory will decide whether to return a client configured for the real Google API or for the new offline service.

- [ ] **Milestone 3: Offline Service Implementation**
  - [ ] Create a new `internal/gtasks/offline_service.go` file.
  - [ ] Implement the `TasksService` interface with an `offlineService` struct.
  - [ ] The methods of `offlineService` will call the corresponding methods on the `offlineStore`.
  - [ ] At this point, existing commands like `gtasks tasks list --offline` should work with the local data.

- [ ] **Milestone 4: The `sync` Command**
  - [ ] Create a new `gtasks sync` command.
  - [ ] Implement the logic to:
    - [ ] Read all data from the `offlineStore`.
    - [ ] Authenticate and connect to the live Google Tasks API.
    - [ ] Implement a basic sync strategy (e.g., push all local changes to the remote).
  - [ ] Add a confirmation prompt before syncing.

- [ ] **Milestone 5: Documentation**
  - [ ] Update the main `README.md` to fully document the `--offline` flag and the new `sync` command.
  - [ ] Update the `DESIGN.md` to reflect the new offline architecture and how it coexists with the online mode.
