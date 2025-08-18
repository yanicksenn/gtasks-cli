# gtasks CLI

A command-line interface (CLI) for managing your Google Tasks.

## 0. Installation

### Homebrew

```sh
brew tap yanicksenn/gtasks-cli
brew install yanicksenn/gtasks-cli
gtasks-cli help
```

### Go

```sh
go install github.com/yanicksenn/gtasks-cli@latest
gtasks-cli help
```

## Configuration

`gtasks-cli` stores its configuration in a file named `config.json` in the following locations:

*   **Linux:** `~/.config/gtasks/config.json`
*   **macOS:** `~/Library/Application Support/gtasks/config.json`
*   **Windows:** `%APPDATA%\gtasks\config.json`

The configuration file has the following structure:

```json
{
  "active_account": "user@example.com"
}
```

*   `active_account`: The email address of the currently active Google account.

## 1. Building and Running

### Prerequisites

*   Go 1.20 or later

### Building

To build the `gtasks` binary, run the following command from the root of the project:

```sh
go build -o gtasks .
```

This will place the executable file named `gtasks` in the project root.

### Running

To run the application, first, you need to authenticate with your Google account:

```sh
./gtasks accounts login
```

This will open a browser window for you to complete the authentication process.

Once authenticated, you can use the other commands, for example:

```sh
./gtasks tasklists list
```

### Running Tests

To run the full suite of tests, use the following command from the root of the project:

```bash
go test ./...
```

This command executes all unit and integration tests against a high-fidelity, in-memory mock of the Google Tasks API, ensuring that no real network calls are made and no authentication is required. It also runs the basic E2E tests.

## Table of Contents

- [1. Authentication](#1-authentication)
- [2. Command Structure](#2-command-structure)
- [3. Offline Mode](#3-offline-mode)
- [4. Terminology](#4-terminology)
- [5. Command Reference](#5-command-reference)
  - [Account Management](#account-management)
  - [TaskList Management](#tasklist-management)
  - [Task Management](#task-management)
- [6. Examples](#6-examples)
- [7. Interactive Mode](#7-interactive-mode)
- [8. Error Handling](#8-error-handling)
- [9. Implementation Details](#9-implementation-details)

---

## 1. Authentication

- **Google Sign-In:** The CLI authenticates with Google using OAuth 2.0.
- **Credential Caching:** Caches credentials locally for automatic use until they expire.
- **Token Refresh:** Automatically refreshes expired tokens.
- **Multi-Account Support:** Manage multiple Google accounts seamlessly.

## 2. Command Structure

The CLI follows a `gtasks <resource> <action> [flags]` pattern.

- **`resource`**: The type of object to operate on (e.g., `accounts`, `tasklists`, `tasks`).
- **`action`**: The operation to perform (e.g., `list`, `create`, `get`, `update`, `delete`).

### Global Flags

- `--offline`: Enable offline mode.
- `--output` (string, optional): Output format. One of `table`, `json`, or `yaml`. Defaults to `table`.
- `--quiet`, `-q` (boolean, optional): Suppress all output.
- `--version`, `-v`: Print the version number.

## 3. Offline Mode

`gtasks` supports a full offline mode. By using the global `--offline` flag, you can manage your tasks and task lists without an internet connection. All changes are saved to a local file (`~/.config/gtasks/offline.json`).

### Syncing Offline Changes

Synchronization must be handled manually. To sync your offline changes, you need to be online and run the commands again without the `--offline` flag. For example, if you created a new task offline:

```sh
./gtasks tasks create --title "My new task" --offline
```

To sync this task with Google Tasks, run the same command again without the `--offline` flag:

```sh
./gtasks tasks create --title "My new task"
```

This will create the task in your Google Tasks. Similarly, for other commands like `update` and `delete`, you need to re-run them without the `--offline` flag to sync the changes.

---

## 4. Terminology

- **Account:** Refers to the Google Account you authenticate with via the SSO sign-in flow. The CLI can cache multiple accounts, but only one is active at a time.
- **TaskList:** A container for your tasks. A user can have multiple task lists to organize different areas of their life (e.g., "Work," "Groceries," "Personal Projects"). Each task list has a unique ID.
- **Task:** A single to-do item that exists within a specific TaskList. It has properties like a title, notes, due date, and a completion status. Each task has a unique ID.
- **@default TaskList:** This is a special identifier that refers to the user's default task list. Google Tasks automatically creates a default list for every user, which is typically named "My Tasks". `gtasks-cli` uses this as the default tasklist for all task-related commands unless a specific tasklist is provided with the `--tasklist` flag.

---

## 5. Command Reference

### Account Management

Manage your authenticated Google accounts.

#### `gtasks accounts login`
Initiates the Google SSO flow to authenticate a new user. The new account becomes the active one.
- **Usage:** `gtasks accounts login`

#### `gtasks accounts logout`
Removes the cached credentials for the currently active user.
- **Usage:** `gtasks accounts logout`

#### `gtasks accounts list`
Lists all authenticated Google accounts.
- **Usage:** `gtasks accounts list`

#### `gtasks accounts switch`
Switches the active user to another authenticated account.
- **Usage:** `gtasks accounts switch <email>`
- **Arguments:**
  - `<email>` (required): The email address of the account to make active.

---

### TaskList Management

Manage your task lists.

#### `gtasks tasklists list`
Lists all task lists.
- **Usage:** `gtasks tasklists list [flags]`
- **Flags:**
  - `--sort-by` (string, optional): Sort task lists by `alphabetical`, `last-modified`, or `uncompleted-tasks`. Defaults to `alphabetical`.

#### `gtasks tasklists get`
Retrieves the details of a specific task list.
- **Usage:** `gtasks tasklists get <tasklist_id>`
- **Arguments:**
  - `<tasklist_id>` (required): The ID of the task list to retrieve.

#### `gtasks tasklists create`
Creates a new task list.
- **Usage:** `gtasks tasklists create --title <list_title>`
- **Flags:**
  - `--title` (string, required): The title for the new task list.

#### `gtasks tasklists update`
Updates the title of an existing task list.
- **Usage:** `gtasks tasklists update <tasklist_id> --title <new_title>`
- **Arguments:**
  - `<tasklist_id>` (required): The ID of the task list to update.
- **Flags:**
  - `--title` (string, required): The new title for the task list.

#### `gtasks tasklists delete`
Permanently deletes a task list and all of its tasks.
- **Usage:** `gtasks tasklists delete <tasklist_id>`
- **Arguments:**
  - `<tasklist_id>` (required): The ID of the task list to delete.



### Task Management

Manage your tasks within a task list.

#### `gtasks tasks list`
Lists tasks within a specific task list.
- **Usage:** `gtasks tasks list [--tasklist <tasklist_id>] [flags]`
- **Flags:**
  - `--tasklist` (string, optional): The ID of the task list. Defaults to `@default`.
  - `--show-completed` (boolean, optional): Include completed tasks.
  - `--show-hidden` (boolean, optional): Include hidden tasks.
  - `--title-contains` (string, optional): Filter tasks by title (case-insensitive).
  - `--notes-contains` (string, optional): Filter tasks by notes (case-insensitive).
  - `--due-before` (string, optional): Filter tasks with a due date before a specified date (e.g., "2025-12-31").
  - `--due-after` (string, optional): Filter tasks with a due date after a specified date (e.g., "2025-12-31").
  - `--sort-by` (string, optional): Sort tasks by `alphabetical`, `last-modified`, or `due-date`. Defaults to `alphabetical`.

#### `gtasks tasks get`
Retrieves the details of a specific task.
- **Usage:** `gtasks tasks get <task_id> [--tasklist <tasklist_id>]`
- **Arguments:**
  - `<task_id>` (required): The ID of the task.
- **Flags:**
  - `--tasklist` (string, optional): The ID of the task list containing the task. Defaults to `@default`.

#### `gtasks tasks create`
Creates a new task in a task list.
- **Usage:** `gtasks tasks create --title <task_title> [--tasklist <tasklist_id>]`
- **Flags:**
  - `--tasklist` (string, optional): The ID of the task list. Defaults to `@default`.
  - `--title` (string, required): The title of the task.
  - `--notes` (string, optional): Notes or description for the task.
  - `--due` (string, optional): Due date in RFC3339 format (e.g., "2025-12-31T22:00:00.000Z").

#### `gtasks tasks update`
Updates an existing task.
- **Usage:** `gtasks tasks update <task_id> [--tasklist <task_id>] [flags]`
- **Arguments:**
  - `<task_id>` (required): The ID of the task to update.
- **Flags:**
  - `--tasklist` (string, optional): The ID of the task list. Defaults to `@default`.
  - `--title` (string, optional): The new title for the task.
  - `--notes` (string, optional): The new notes for the task.
  - `--due` (string, optional): The new due date in RFC3339 format.

#### `gtasks tasks complete`
Marks a task as complete.
- **Usage:** `gtasks tasks complete <task_id> [--tasklist <tasklist_id>]`
- **Arguments:**
  - `<task_id>` (required): The ID of the task.
- **Flags:**
  - `--tasklist` (string, optional): The ID of the task list. Defaults to `@default`.

#### `gtasks tasks uncomplete`
Marks a task as not complete.
- **Usage:** `gtasks tasks uncomplete <task_id> [--tasklist <tasklist_id>]`
- **Arguments:**
  - `<task_id>` (required): The ID of the task.
- **Flags:**
  - `--tasklist` (string, optional): The ID of the task list. Defaults to `@default`.

#### `gtasks tasks delete`
Permanently deletes a task.
- **Usage:** `gtasks tasks delete <task_id> [--tasklist <tasklist_id>]`
- **Arguments:**
  - `<task_id>` (required): The ID of the task.
- **Flags:**
  - `--tasklist` (string, optional): The ID of the task list. Defaults to `@default`.



## 6. Examples

Here are some common commands to get you started.

### Authenticate with Google
This is the first command you should run. It will open your web browser to the Google login page to authorize the application.
```sh
./gtasks accounts login
```

### List All Your Task Lists
```sh
$ ./gtasks tasklists list
Task Lists:
- My tasks (MTM1NTM2MzQzNzczNDkyNzc1NTQ6MDow)
- Old Google Keep reminders (eUZqdFdsOGpsNVdUclY1Mg)
```

### Create a New Task List
```sh
$ ./gtasks tasklists create --title "Groceries"
Successfully created task list: Groceries (OS0ydmR2N3NpSTQ4SzVVMA)
```

### List Tasks
You can list tasks in your default list or a specific list.

```sh
# List tasks in the default list
$ ./gtasks tasks list
Tasks:
[ ] buy toilet cleaners (U0QzVTI3TDFiRXg1NnJoSg)
[ ] buy shampoo (VlFTcEt1TXItMl9RUDZpRg)
...

# List tasks in a specific list by its ID
$ ./gtasks tasks list --tasklist "OS0ydmR2N3NpSTQ4SzVVMA"
Tasks:
[ ] Buy milk (ZmFyb3FBSzJhUUlRZGJnWg)

# Find all tasks with "buy" in the title
$ gtasks tasks list --title-contains "buy"
Tasks:
[ ] buy toilet cleaners (U0QzVTI3TDFiRXg1NnJoSg)
[ ] buy shampoo (VlFTcEt1TXItMl9RUDZpRg)
```

### Create a New Task
You can create a simple task or add more details like a due date.

```sh
# Create a simple task in the default list
$ ./gtasks tasks create --title "Buy milk"
Successfully created task: Buy milk (eG9_b...)

# Create a task with a due date in a specific list
$ ./gtasks tasks create --tasklist "OS0ydmR2N3NpSTQ4SzVVMA" --title "Finish report" --due "2025-12-20T15:00:00.000Z"
Successfully created task: Finish report (aG9_c...)
```

### Complete a Task
To complete a task, you need its ID, which you can get from the `tasks list` command.
```sh
$ ./gtasks tasks complete "ZmFyb3FBSzJhUUlRZGJnWg" --tasklist "OS0ydmR2N3NpSTQ4SzVVMA"
Successfully completed task: Buy milk (ZmFyb3FBSzJhUUlRZGJnWg)
```

### View Completed Tasks
Use the `--show-completed` flag to include completed tasks in the list.
```sh
$ ./gtasks tasks list --tasklist "OS0ydmR2N3NpSTQ4SzVVMA" --show-completed
Tasks:
[x] Buy milk (ZmFyb3FBSzJhUUlRZGJnWg)
```

### Work Offline
You can use the `--offline` flag with most commands to work with a local copy of your tasks.
```sh
./gtasks tasklists list --offline
```

### Advanced Examples

#### Filtering Tasks

You can combine filters to narrow down your search.

```sh
# Find tasks with "report" in the title that are due before the end of 2025
$ ./gtasks tasks list --title-contains "report" --due-before "2025-12-31"

# Find tasks with "meeting" in the notes
$ ./gtasks tasks list --notes-contains "meeting"
```

#### Changing the Output Format

You can change the output format to JSON or YAML, which is useful for scripting.

```sh
# Get a list of task lists in JSON format
$ ./gtasks tasklists list --output json
[
  {
    "id": "MTM1NTM2MzQzNzczNDkyNzc1NTQ6MDow",
    "title": "My tasks",
    ...
  }
]
```

## 7. Interactive Mode

`gtasks` provides a full-screen interactive mode that allows you to manage your tasks in a more fluid, application-like experience.

To start the interactive mode, run the following command:

```sh
./gtasks interactive --tasklist <tasklist_id>
```

### Features

*   **View Tasks:** See all your tasks in a list.
*   **Add Tasks:** Press `a` to add a new task.
*   **Edit Tasks:** Press `e` to edit the selected task.
*   **Delete Tasks:** Press `d` to delete the selected task.
*   **Complete/Uncomplete Tasks:** Press `space` to toggle the completion status of the selected task.

### Keybindings

-   `q`, `ctrl+c`: Quit the application.
-   `up`, `k`: Move the cursor up.
-   `down`, `j`: Move the cursor down.
-   `a`: Add a new task.
-   `e`: Edit the selected task.
-   `d`: Delete the selected task.
-   `space`: Toggle the completion status of the selected task.

## 8. Error Handling

`gtasks-cli` is designed to provide clear and actionable feedback when errors occur. Here's how it handles common issues:

*   **Network Errors:** If the CLI cannot connect to the Google Tasks API, it will print an error message and exit. If you are offline, you can use the `--offline` flag to work with a local copy of your tasks.
*   **Authentication Errors:** If your credentials have expired or are invalid, the CLI will prompt you to log in again.
*   **API Errors:** If the Google Tasks API returns an error, the CLI will print the error message from the API and exit.
*   **Not Found Errors:** If you try to access a resource that does not exist (e.g., a task or tasklist with an invalid ID), the CLI will print a "not found" error and exit.

## 9. Implementation Details

- **Language:** Go
- **Libraries:**
  - Cobra (`github.com/spf13/cobra`) for CLI structure.
  - Bubble Tea (`github.com/charmbracelet/bubbletea`) for the interactive TUI.
  - Google API Client for Go (`google.golang.org/api/tasks/v1`).
  - Go OAuth2 Library (`golang.org/x/oauth2`).