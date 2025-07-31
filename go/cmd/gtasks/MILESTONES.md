# Milestones: Authentication Flow Overhaul

This document outlines the high-level milestones for transitioning the `gtasks-cli` from a local OAuth flow to a centralized, web-based authentication system.

---

### Milestone 1: Implement the Core Web Authentication Flow

*Goal: Create the foundational client-side logic for the new authentication method.*

-   [x] **Create `internal/auth/web.go`:** This new file will contain the primary function for orchestrating the web-based login process.
-   [x] **Add Dependencies:** Successfully add and verify the `github.com/pkg/browser` package.
-   [x] **Implement Local Callback Server:** The CLI must be able to start a local server to listen for the redirect from the external authentication service.
-   [x] **Implement Browser Handling:** The CLI must be able to reliably open the user's default web browser to the correct authentication URL.

---

### Milestone 2: Integrate and Refactor the CLI

*Goal: Fully replace the old authentication flow with the new one and clean up obsolete code.*

-   [x] **Modify `login` Command:** Update the `accounts login` command to exclusively use the new `LoginViaWebFlow` function.
-   [x] **Remove Obsolete Code:** Delete the old local OAuth logic from `internal/auth/auth.go`, including the `NewClient` and `NewAuthenticator` functions.
-   [x] **Verify Token Caching:** Ensure that the new flow correctly saves tokens to the existing cache and that user switching (`accounts switch`) remains functional.

---

### Milestone 3: Testing and Documentation

*Goal: Ensure the new implementation is robust, reliable, and well-documented.*

-   [x] **Update Unit Tests:** Adapt existing tests for the `login` command to work with the new asynchronous, web-based flow, likely using a mock HTTP server.
-   [ ] **End-to-End Manual Testing:** Perform thorough manual testing of the complete login, logout, and account switching process.
-   [ ] **Update `README.md`:** Revise the main project `README.md` to reflect the new, simplified login procedure for end-users.
-   [ ] **Update `DESIGN.md`:** Update the internal design documentation to describe the new authentication architecture.
