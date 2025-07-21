# gtasks CLI - Implementation Milestones

This document outlines the development plan as a series of ordered milestones. We will track our progress by checking off each item as it's completed and tested.

- [x] **Milestone 1: Foundation & Scaffolding**
  - [x] Create the directory structure as defined in `DESIGN.md`.
  - [x] Initialize the Go module (`go mod init`).
  - [x] Add initial dependencies (`cobra`, `oauth2`, `google-api-go-client`).
  - [x] Set up the main application entry point (`main.go`) and the root Cobra command.

- [x] **Milestone 2: Core Authentication**
  - [x] Implement the `gtasks login` command.
  - [x] Implement the full OAuth2 web flow to retrieve a token.
  - [x] Implement secure token storage in the user's home directory.
  - [x] Implement the configuration logic to track the active user.
  - [x] Create offline tests for the auth and config logic.

- [x] **Milestone 3: TaskList Read Operations**
  - [x] Implement the `gtasks tasklists list` command.
  - [x] Implement the `gtasks tasklists get` command.
  - [x] Create the internal client wrapper for the Google Tasks API.
  - [x] Create offline tests using mock interfaces for the TaskList read operations.

- [x] **Milestone 4: TaskList Write Operations**
  - [x] Implement the `gtasks tasklists create` command.
  - [x] Implement the `gtasks tasklists update` command.
  - [x] Implement the `gtasks tasklists delete` command.
  - [x] Add offline tests for the TaskList write operations.

- [ ] **Milestone 5: Task Read Operations**
  - [ ] Implement the `gtasks tasks list` command.
  - [ ] Implement the `gtasks tasks get` command.
  - [ ] Add offline tests for the Task read operations.

- [ ] **Milestone 6: Task Write Operations**
  - [ ] Implement the `gtasks tasks create` command.
  - [ ] Implement the `gtasks tasks update` command.
  - [ ] Implement the `gtasks tasks complete` command.
  - [ ] Implement the `gtasks tasks delete` command.
  - [ ] Add offline tests for the Task write operations.

- [ ] **Milestone 7: Account Management**
  - [ ] Implement the `gtasks logout` command.
  - [ ] Implement the `gtasks accounts list` command.
  - [ ] Implement the `gtasks accounts switch` command.
  - [ ] Add offline tests for the account management logic.

- [ ] **Milestone 8: End-to-End Testing**
  - [ ] Set up the E2E test suite structure.
  - [ ] Write E2E tests covering the main user flows for `tasklists` and `tasks`.
  - [ ] Document how to run the E2E tests.
