# TODO: Refactoring the `cmd` Package for Clarity and Reusability

## Problem

The command definitions located in the `cmd/` package (specifically `tasks.go`, `tasklists.go`, and `accounts.go`) contain a significant amount of repetitive boilerplate code within their `RunE` functions. This pattern is repeated across most commands:

1.  A `gtasks.Client` is created using `gtasks.NewClient(cmd, context.Background())`.
2.  The error from the client creation is handled with `fmt.Errorf("error creating client: %w", err)`.
3.  The `--quiet` flag is checked with `cmd.Flags().GetBool("quiet")` to determine if output should be suppressed.
4.  Configuration is loaded with `config.Load()` and its associated error handling.

This leads to verbose and cluttered `RunE` functions where the core command logic is obscured by setup code.

## Why It's Bad

This repetition of boilerplate code is detrimental to the project for several reasons:

-   **Violation of DRY Principle (Don't Repeat Yourself):** It introduces redundancy, making the codebase larger and more difficult to maintain than necessary.
-   **Reduced Readability:** The essential logic of each command is buried within layers of setup code, making it hard to quickly understand a command's purpose and implementation.
-   **Maintenance Challenges:** If the client or config initialization logic ever needs to be changed (e.g., to add a new parameter or modify error handling), the change must be manually applied to every single command. This is tedious, inefficient, and highly susceptible to human error.

## Proposed Solution

We will refactor the `cmd` package to eliminate this boilerplate by abstracting the common setup and teardown logic.

1.  **Introduce a `CommandHelper` struct.**
    *   Create a new helper struct (e.g., in `cmd/helpers.go`) that encapsulates the common objects required by commands.
    *   This struct will hold the initialized `gtasks.Client`, the loaded `*config.Config`, and the values of common persistent flags like `--quiet` and `--offline`.

2.  **Create a factory function for the helper.**
    *   Develop a constructor function, such as `NewCommandHelper(cmd *cobra.Command)`, that performs the repetitive setup tasks once:
        *   It will load the configuration.
        *   It will initialize the appropriate `gtasks.Client` (online or offline).
        *   It will parse the persistent flags.
        *   It will return an initialized `CommandHelper` instance or an error.

3.  **Refactor all commands to use the helper.**
    *   Modify the `RunE` function of each command in `tasks.go`, `tasklists.go`, and `accounts.go` to use the `CommandHelper`.
    *   The body of each `RunE` will be reduced to two main parts:
        1.  A single call to `NewCommandHelper(cmd)`.
        2.  The core logic of the command, which now uses the fields from the helper (e.g., `helper.Client`, `helper.Config`, `helper.Quiet`).

This approach will dramatically simplify the command definitions, making them cleaner, more readable, and much easier to maintain. An example of the refactored command would look like this:

```go
// Example of a refactored command
var listTasksCmd = &cobra.Command{
    Use:   "list",
    Short: "List all tasks in a task list",
    RunE: func(cmd *cobra.Command, args []string) error {
        helper, err := NewCommandHelper(cmd)
        if err != nil {
            return err // Centralized error handling
        }

        // Core logic is now clean and focused
        tasks, err := helper.Client.ListTasks(...)
        if err != nil {
            return err
        }

        if !helper.Quiet {
            helper.UI.PrintTasks(tasks) // Using a dedicated UI printer
        }
        return nil
    },
}
```
