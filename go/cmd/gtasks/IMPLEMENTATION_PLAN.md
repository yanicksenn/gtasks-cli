# gtasks CLI - Implementation Plan

**Related Documents:**
- [Software Design (`DESIGN.md`)](./DESIGN.md)

---

This document lists the expected files and data structures required to build the `gtasks` CLI, based on the `DESIGN.md`.

## 1. File Structure

This is the planned layout of the Go source files within the `go/cmd/gtasks/` directory.

```
.
├── DESIGN.md
├── IMPLEMENTATION_PLAN.md
├── README.md
├── cmd/
│   ├── accounts.go      # `gtasks accounts` commands (list, switch)
│   ├── root.go          # Root `gtasks` command and global flags
│   ├── tasklists.go     # `gtasks tasklists` commands (list, create, get, etc.)
│   └── tasks.go         # `gtasks tasks` commands (list, create, get, etc.)
├── internal/
│   ├── auth/
│   │   └── auth.go      # Handles OAuth2 flow, token storage, and client creation.
│   ├── config/
│   │   └── config.go    # Manages user configuration (e.g., active account).
│   └── gtasks/
│       ├── client.go    # The main client wrapper for the Google Tasks service.
│       ├── tasklists.go # Business logic for tasklist operations.
│       └── tasks.go     # Business logic for task operations.
└── main.go              # Main application entry point.
```

## 2. Key Data Structures

These are the primary structs that will be used to pass data between the CLI and the business logic layers.

### Configuration & Authentication Structs

Located in `internal/config/config.go`:
```go
// Config represents the application's configuration.
type Config struct {
    ActiveAccount string `json:"active_account"`
}
```

Located in `internal/auth/auth.go`:
```go
// TokenCache represents the structure of the credentials file.
// It maps a user's email to their OAuth2 token.
type TokenCache struct {
    Tokens map[string]*oauth2.Token `json:"tokens"`
}
```

### TaskList Management Structs

Located in `internal/gtasks/tasklists.go`.

#### **List TaskLists**
- **API Documentation:** [`tasklists.list`](https://developers.google.com/tasks/reference/rest/v1/tasklists/list)
- No options struct needed for this operation.

#### **Get TaskList**
- **API Documentation:** [`tasklists.get`](https://developers.google.com/tasks/reference/rest/v1/tasklists/get)
```go
type GetTaskListOptions struct {
    TaskListID string
}
```

#### **Create TaskList**
- **API Documentation:** [`tasklists.insert`](https://developers.google.com/tasks/reference/rest/v1/tasklists/insert)
```go
type CreateTaskListOptions struct {
    Title string
}
```

#### **Update TaskList**
- **API Documentation:** [`tasklists.update`](https://developers.google.com/tasks/reference/rest/v1/tasklists/update)
```go
type UpdateTaskListOptions struct {
    TaskListID string
    Title      string
}
```

#### **Delete TaskList**
- **API Documentation:** [`tasklists.delete`](https://developers.google.com/tasks/reference/rest/v1/tasklists/delete)
```go
type DeleteTaskListOptions struct {
    TaskListID string
}
```

### Task Management Structs

Located in `internal/gtasks/tasks.go`.

#### **List Tasks**
- **API Documentation:** [`tasks.list`](https://developers.google.com/tasks/reference/rest/v1/tasks/list)
```go
type ListTasksOptions struct {
    TaskListID     string
    ShowCompleted  bool
    ShowHidden     bool
}
```

#### **Get Task**
- **API Documentation:** [`tasks.get`](https://developers.google.com/tasks/reference/rest/v1/tasks/get)
```go
type GetTaskOptions struct {
    TaskListID string
    TaskID     string
}
```

#### **Create Task**
- **API Documentation:** [`tasks.insert`](https://developers.google.com/tasks/reference/rest/v1/tasks/insert)
```go
type CreateTaskOptions struct {
    TaskListID string
    Title      string
    Notes      string
    Due        string // Using string for simplicity, will be parsed to RFC3339
}
```

#### **Update Task**
- **API Documentation:** [`tasks.update`](https://developers.google.com/tasks/reference/rest/v1/tasks/update)
```go
type UpdateTaskOptions struct {
    TaskListID string
    TaskID     string
    Title      string
    Notes      string
    Due        string
}
```

#### **Complete Task**
- **Note:** This is an `update` operation that sets the task `status` to `"completed"`.
- **API Documentation:** [`tasks.update`](https://developers.google.com/tasks/reference/rest/v1/tasks/update)
```go
type CompleteTaskOptions struct {
    TaskListID string
    TaskID     string
}
```

#### **Delete Task**
- **API Documentation:** [`tasks.delete`](https://developers.google.com/tasks/reference/rest/v1/tasks/delete)
```go
type DeleteTaskOptions struct {
    TaskListID string
    TaskID     string
}
```
