# Refactoring Milestones

This document tracks the progress of the codebase refactoring. Each major step is a milestone that will be checked off upon completion.

- [x] **1. Project Scaffolding:** Create the necessary directory structure for the refactoring.
- [x] **2. Unify Data Stores:** Merge the `offline_store` and `mock_store` into a single, reusable `in_memory_store`.
- [ ] **3. Separate UI Concerns:** Move all user-facing printing logic into a dedicated `internal/ui` package.
- [ ] **4. Refactor Command Boilerplate:** Abstract repetitive setup code in the `cmd` package into a `CommandHelper`.
- [ ] **5. Externalize Credentials:** Move hardcoded OAuth2 credentials out of the source code.
- [ ] **6. Final Cleanup:** Remove temporary `TODO` and `MILESTONES` documents.
