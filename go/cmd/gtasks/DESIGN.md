# gtasks CLI - Software Design

**Related Documents:**
- [Requirements (`README.md`)](./README.md)
- [Implementation Plan (`IMPLEMENTATION_PLAN.md`)](./IMPLEMENTATION_PLAN.md)

---

This document outlines the high-level software design for the `gtasks` CLI tool.

## 1. Overall Architecture

The application will be built in Go and structured around the **Cobra** library (`github.com/spf13/cobra`). Cobra provides a robust framework for creating modern, hierarchical CLI applications, making it easy to manage commands, subcommands, and flags.

The core structure will be:
- A `root` command (`gtasks`) that serves as the entry point.
- Resource-based subcommands (`accounts`, `tasklists`, `tasks`).
- Action-based sub-subcommands (`list`, `create`, `get`, etc.).

This approach ensures the CLI is organized, discoverable, and easy to extend.

## 2. Directory Structure

The project will follow standard Go conventions, organized for clarity and maintainability within the `go/cmd/gtasks/` directory:

```
gtasks/
├── cmd/          # Cobra command definitions
│   ├── root.go
│   ├── accounts.go
│   ├── tasklists.go
│   └── tasks.go
├── internal/     # Internal application logic, not for export
│   ├── auth/     # Google OAuth2 authentication and token management
│   ├── config/   # Configuration handling (active user, credentials path)
│   └── gtasks/   # Wrapper for the Google Tasks API client
└── main.go       # Main application entry point
```

## 3. Authentication Flow (OAuth 2.0)

Authentication is the most critical component. The flow will be:

1.  **Trigger:** When a command requiring authentication is run for the first time (or after `gtasks login`), the `auth` module is invoked.
2.  **OAuth Config:** The application will use a pre-configured OAuth 2.0 Client ID and Secret for a "Desktop" application, obtained from the Google Cloud Console.
3.  **Consent URL:** The user will be presented with a URL to open in their browser.
4.  **Authorization Code:** After the user grants consent, Google will redirect them to a local callback URL (e.g., `http://localhost:8080`). The CLI will run a temporary local web server to capture the authorization code from this redirect.
5.  **Token Exchange:** The application will exchange the authorization code for an access token and a refresh token.
6.  **Credential Storage:** The tokens (especially the refresh token) will be stored in a JSON file located in the user's home directory (e.g., `~/.config/gtasks/credentials.json`).

The refresh token will be used to automatically get new access tokens when the old ones expire, providing a seamless experience after the initial login.

## 4. Configuration Management

A simple configuration file (e.g., `~/.config/gtasks/config.yaml`) will manage application settings. This will primarily store:
- The email/ID of the currently **active account**. This allows the CLI to know which set of credentials to use from the `credentials.json` file when multiple accounts are authenticated.

The `config` package will handle reading and writing this configuration.

## 5. Google Tasks API Integration

A dedicated `gtasks` package within `internal/` will act as a wrapper around the official `google.golang.org/api/tasks/v1` client.

- **Initialization:** This package will be responsible for creating an authenticated `http.Client` using the OAuth2 tokens managed by the `auth` module.
- **Service Creation:** It will use the authenticated client to create an instance of the `tasks.Service`.
- **API Abstraction:** It will provide high-level functions that map directly to the CLI actions (e.g., `ListTaskLists()`, `CreateTask(opts CreateTaskOptions)`). This decouples the Cobra command logic from the direct Google API calls.

### Command Structs for Decoupling

To ensure a clean separation of concerns, we will use dedicated structs to pass data from the CLI layer to the business logic layer.

1.  **CLI Layer (`cmd/`):** The Cobra command's `Run` function will be responsible *only* for parsing flags and arguments.
2.  **Command Options Struct:** It will populate a specific options struct (e.g., `gtasks.CreateTaskOptions`) with the parsed data.
3.  **Business Logic Layer (`internal/gtasks/`):** This struct is then passed to the corresponding function in the `gtasks` package (e.g., `gtasks.CreateTask(opts)`).
4.  **Execution:** The business logic function unpacks the struct and performs the necessary Google Tasks API calls.

This pattern makes the `gtasks` package completely independent of Cobra, allowing it to be tested in isolation and potentially reused by other interfaces in the future.

## 6. Testing Strategy

The project will employ a two-tiered testing strategy to ensure correctness and reliability.

### Offline Tests (Unit/Integration)

- **Goal:** To verify the internal business logic without making any real network calls.
- **Location:** Tests will be co-located with the code they test (e.g., `internal/gtasks/tasklists_test.go`).
- **Method:** We will define interfaces for our Google Tasks API interactions. In tests, these interfaces will be fulfilled by mock implementations that simulate the behavior of the Google Tasks API. This allows us to test all logic, including request parameter construction, response handling, and error conditions, in a fast and deterministic manner.

### End-to-End (E2E) Tests

- **Goal:** To verify that the compiled CLI application functions correctly from a user's perspective.
- **Location:** These tests will reside in a separate top-level directory (e.g., `e2e/`).
- **Method:** The E2E test suite will compile the `gtasks` binary and execute it as a subprocess. Tests will invoke actual commands (e.g., `gtasks tasklists list --output json`) and assert against the stdout, stderr, and exit codes. These tests will require a live, authenticated Google account and will be run against the actual Google Tasks API. They are slower and will be run more selectively, but provide the highest level of confidence.