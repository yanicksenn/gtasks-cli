# gtasks CLI - Implementation Milestones

This document outlines the development plan as a series of ordered milestones. We will track our progress by checking off each item as it's completed and tested.

## Table of Contents

- [Milestone 1: Foundation & Scaffolding](#milestone-1-foundation--scaffolding)
- [Milestone 2: Core Authentication](#milestone-2-core-authentication)
- [Milestone 3: TaskList Read Operations](#milestone-3-tasklist-read-operations)
- [Milestone 4: TaskList Write Operations](#milestone-4-tasklist-write-operations)
- [Milestone 5: Task Read Operations](#milestone-5-task-read-operations)
- [Milestone 6: Task Write Operations](#milestone-6-task-write-operations)
- [Milestone 7: Account Management](#milestone-7-account-management)
- [Milestone 8: End-to-End Testing](#milestone-8-end-to-end-testing)

---

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

- [x] **Milestone 5: Task Read Operations**
  - [x] Implement the `gtasks tasks list` command.
  - [x] Implement the `gtasks tasks get` command.
  - [x] Add offline tests for the Task read operations.

- [x] **Milestone 6: Task Write Operations**
  - [x] Implement the `gtasks tasks create` command.
  - [x] Implement the `gtasks tasks update` command.
  - [x] Implement the `gtasks tasks complete` command.
  - [x] Implement the `gtasks tasks delete` command.
  - [x] Add offline tests for the Task write operations.

- [x] **Milestone 7: Account Management**
  - [x] Implement the `gtasks logout` command.
  - [x] Implement the `gtasks accounts list` command.
  - [x] Implement the `gtasks accounts switch` command.
  - [x] Add offline tests for the account management logic.

- [x] **Milestone 8: End-to-End Testing**
  - [x] Set up the E2E test suite structure.
  - [x] Write E2E tests covering the main user flows for `tasklists` and `tasks`.
  - [x] Document how to run the E2E tests.
