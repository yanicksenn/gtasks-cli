# Data Storage

The `gtasks-cli` tool stores data in two main locations:

*   **Configuration:** `~/.config/gtasks/config.json`
*   **Offline Data:** `~/.config/gtasks/offline.json`

## Configuration File

The configuration file (`config.json`) stores the user's settings. It has the following structure:

```json
{
  "active_account": "user@example.com"
}
```

*   `active_account`: The email address of the currently active Google account.

## Offline Data File

The offline data file (`offline.json`) stores a local copy of the user's tasks and task lists. This allows the user to work offline and sync their changes later. The file has the following structure:

```json
{
  "task_lists": {
    "taskListId1": {
      "id": "taskListId1",
      "title": "Task List 1"
    },
    "taskListId2": {
      "id": "taskListId2",
      "title": "Task List 2"
    }
  },
  "tasks": {
    "taskListId1": {
      "taskId1": {
        "id": "taskId1",
        "title": "Task 1",
        "notes": "Notes for Task 1",
        "due": "2025-12-31T22:00:00.000Z",
        "status": "needsAction"
      }
    }
  },
  "next_id": 3
}
```

*   `task_lists`: A map of task lists, where the key is the task list ID.
*   `tasks`: A map of tasks, where the key is the task list ID, and the value is a map of tasks, where the key is the task ID.
*   `next_id`: The next available ID for a new task or task list.
