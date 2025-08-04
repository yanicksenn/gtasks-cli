# Feature Milestone: Advanced Output Formatting

This document tracks the implementation of the advanced output formatting feature.

- [x] **1. Core `ui.Printer` Implementation**
    - [ ] Create the `ui.Printer` struct in the `internal/ui` package.
    - [ ] Add a persistent `--output` flag to `cmd/root.go`.
    - [ ] Integrate the `ui.Printer` into the `cmd.CommandHelper`.

- [x] **2. Refactor `tasklists` Commands**
    - [ ] Refactor `tasklists list` to use the `ui.Printer`.
    - [ ] Refactor `tasklists get` to use the `ui.Printer`.
    - [ ] Remove the `tasklists print` command.

- [x] **3. Refactor `tasks` Commands**
    - [ ] Refactor `tasks list` to use the `ui.Printer`.
    - [ ] Refactor `tasks get` to use the `ui.Printer`.
    - [ ] Remove the `tasks print` command.

- [x] **4. Refactor `accounts` Command**
    - [ ] Implement `Printer.PrintAccounts([]string)`.
    - [ ] Refactor `accounts list` to use the `ui.Printer`.

- [x] **5. Update Tests**
    - [ ] Add unit tests for the `ui.Printer` for all data types and formats.
    - [ ] Update E2E tests to validate the `--output` flag.
    - [ ] Remove E2E tests for the deleted `print` commands.

- [ ] **6. Final Cleanup**
    - [ ] Remove this `MILESTONES.md` document.
    - [ ] Remove the `TODO_FEATURE_OUTPUT.md` document.
