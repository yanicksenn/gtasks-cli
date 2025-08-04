# TODO: Feature - Interactive TUI Mode

## 1. Goals

The primary goal of this feature is to enhance the user experience by providing a rich, interactive, full-screen terminal UI (TUI). This moves the application from a simple command-and-response model to a more fluid, application-like experience within the terminal.

-   **High-Level Goals:**
    -   Create a new command, `gtasks interactive`, to launch the TUI.
    -   Allow users to navigate their task lists and tasks using keyboard shortcuts (e.g., arrow keys, `j`/`k`).
    -   Enable users to perform all major CRUD (Create, Read, Update, Delete) operations on tasks and task lists without leaving the interface.
    -   Improve feature discoverability by presenting available actions and context-sensitive help within the UI.

-   **Key UI Components:**
    -   A multi-pane layout, likely with one pane for task lists and another for the tasks within the selected list.
    -   A status bar to display the current user, connection status (online/offline), and keybindings.
    -   Forms for creating and editing tasks and task lists.
    -   Confirmation dialogs for destructive actions like deletion.

## 2. Technical Trade-offs

The main technical decision is the choice of Go library to build the TUI.

-   **Option A: Bubble Tea (`charmbracelet/bubbletea`)**
    -   **Description:** A modern, functional framework based on The Elm Architecture (TEA). It's highly extensible and state-driven.
    -   **Pros:**
        -   **Robust State Management:** TEA makes application state predictable and easy to manage, which is crucial for a complex UI.
        -   **Flexible & Composable:** Excellent for building custom, bespoke layouts.
        -   **Great Ecosystem:** Part of the Charm ecosystem, with well-integrated libraries for styling (`lipgloss`), components (`bubbles`), and more.
    -   **Cons:**
        -   **Learning Curve:** The TEA pattern can be unfamiliar to developers accustomed to traditional imperative or object-oriented UI frameworks.
        -   **Boilerplate:** Can require more initial setup code compared to widget-based libraries.

-   **Option B: tview (`rivo/tview`)**
    -   **Description:** A more traditional, widget-based TUI library. It provides a rich set of pre-built components like lists, forms, and modals.
    -   **Pros:**
        -   **Rapid Prototyping:** The extensive set of built-in widgets can make it faster to get a functional UI up and running.
        -   **Familiar Model:** Its imperative, component-based approach is more familiar to many developers.
    -   **Cons:**
        -   **State Management:** Can become difficult to manage application state as the UI grows in complexity. State is often spread across different widgets.
        -   **Less Flexible:** Can be more restrictive when you need to create highly customized UI components or layouts.

-   **Recommendation:**
    **Bubble Tea** is the recommended choice. While it has a slightly steeper learning curve, its superior state management model is a major advantage for a stateful application like this. It will lead to a more maintainable and scalable codebase in the long run.

## 3. Testability

Testing TUIs is inherently challenging because you cannot easily assert what is visually rendered to the screen. The testing strategy must therefore focus on the application's logic and state, not its visual output.

-   **Unit Testing the Model (Bubble Tea):**
    -   The core of a Bubble Tea application is the `Update` function, which is a pure function: `(msg, model) -> (model, cmd)`.
    -   We can write unit tests that send specific messages (e.g., a key press, a data-loaded message) to the `Update` function with a given initial `model` (state).
    -   We can then assert that the returned `model` has been updated correctly and that the correct `cmd` (command to be executed, like an API call) is issued. This allows us to test the application's entire logical flow without ever needing to render the UI.

-   **Separation of Concerns:**
    -   The business logic (e.g., calling the Google Tasks API) will be triggered via commands. This logic already exists in our `gtasks` package and is independently testable.
    -   The TUI code will be responsible for dispatching these commands and handling the results (as new messages), keeping the UI and business logic cleanly separated.

-   **E2E Testing:**
    -   Automated E2E testing of a TUI is complex and would require specialized tools for terminal session recording and playback.
    -   For this project, we will rely on thorough unit testing of the TUI's state logic and manual testing for validating the final visual output and user experience.
