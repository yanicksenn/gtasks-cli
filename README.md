# gtasks CLI

A command-line interface (CLI) for managing your Google Tasks.

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
- [8. Implementation Details](#8-implementation-details)
- [9. Project Documentation](#9-project-documentation)
- [10. Running Tests](#10-running-tests)
- [11. Homebrew Release](#11-homebrew-release)

---

## Building and Running

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

**Note:** Synchronization must be handled manually. This tool does not currently provide a `sync` command.

---

## 4. Terminology

- **Account:** Refers to the Google Account you authenticate with via the SSO sign-in flow. The CLI can cache multiple accounts, but only one is active at a time.
- **TaskList:** A container for your tasks. A user can have multiple task lists to organize different areas of their life (e.g., "Work," "Groceries," "Personal Projects"). Each task list has a unique ID.
- **Task:** A single to-do item that exists within a specific TaskList. It has properties like a title, notes, due date, and a completion status. Each task has a unique ID.

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
- **Usage:** `gtasks tasks update <task_id> [--tasklist <tasklist_id>] [flags]`
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

### Keybindings

-   `q`, `ctrl+c`: Quit the application.
-   `up`, `k`: Move the cursor up.
-   `down`, `j`: Move the cursor down.
-   `space`: Toggle the completion status of the selected task.

## 8. Implementation Details

- **Language:** Go
- **Libraries:**
  - Cobra (`github.com/spf13/cobra`) for CLI structure.
  - Bubble Tea (`github.com/charmbracelet/bubbletea`) for the interactive TUI.
  - Google API Client for Go (`google.golang.org/api/tasks/v1`).
  - Go OAuth2 Library (`golang.org/x/oauth2`).

## 9. Project Documentation

For more detailed information on the design and implementation, see the documents in the following directory:

- [Software Design (`/design`)](./design)

## 10. Running Tests

To run the full suite of tests, use the following command from the root of the project:

```bash
go test ./...
```

This command executes all unit and integration tests against a high-fidelity, in-memory mock of the Google Tasks API, ensuring that no real network calls are made and no authentication is required. It also runs the basic E2E tests.

## 11. Homebrew Release

To release a new version of `gtasks-cli` to Homebrew, follow these steps:

### 1. Create a Git Tag

First, create a new Git tag for the release. This tag will be used to identify the specific version of the code to be released.

```sh
git tag v0.1.0
git push origin v0.1.0
```

### 2. Create a Source Archive

Next, create a source archive (tarball) of the tagged release. This archive will be uploaded to GitHub and used by Homebrew to download the source code.

```sh
git archive --format=tar.gz --output=gtasks-cli-v0.1.0.tar.gz v0.1.0
```

### 3. Calculate the SHA256 Checksum

Homebrew uses a SHA256 checksum to verify the integrity of the downloaded source archive. Calculate the checksum of the tarball you just created.

```sh
shasum -a 256 gtasks-cli-v0.1.0.tar.gz
```

### 4. Create or Update the Homebrew Formula

A Homebrew formula is a Ruby file that tells Homebrew how to install your package. You will need to create a new formula file or update an existing one.

The formula file should be named `gtasks-cli.rb` and should look like this:

```ruby
class GtasksCli < Formula
  desc "A CLI for managing Google Tasks"
  homepage "https://github.com/yanicksenn/gtasks-cli"
  url "https://github.com/yanicksenn/gtasks-cli/archive/v0.1.0.tar.gz"
  sha256 "YOUR_SHA256_CHECKSUM_HERE"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "."
  end

  test do
    system "#{bin}/gtasks", "--help"
  end
end
```

Replace `YOUR_SHA256_CHECKSUM_HERE` with the checksum you calculated in the previous step.

### 5. Create a Homebrew Tap

A Homebrew Tap is a Git repository that contains your formula files. If you don't already have one, you'll need to create a new public GitHub repository named `homebrew-tap`.

### 6. Add the Formula to Your Tap

Once you have a Tap, create a `Formula` directory inside it and move your `gtasks-cli.rb` file into that directory.

```sh
mkdir -p homebrew-tap/Formula
mv gtasks-cli.rb homebrew-tap/Formula/
```

Commit and push the changes to your `homebrew-tap` repository.

### 7. Install the CLI

Users can now install `gtasks-cli` by first tapping your repository and then running the install command:

```sh
brew tap yanicksenn/tap
brew install gtasks-cli
```