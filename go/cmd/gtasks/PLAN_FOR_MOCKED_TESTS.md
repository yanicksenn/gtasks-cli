# Plan for Mocked Tests

This document outlines the plan for implementing a mock Google Tasks API server to enable robust, offline testing of our CLI.

## 1. The Goal

The primary goal is to create a test environment where our application's logic can be verified without making any real network calls to Google. This will be achieved by running a local, in-memory mock server that simulates the behavior of the Google Tasks API.

## 2. The Approach

We will use Go's standard library, primarily the `net/http/httptest` package, to create a mock HTTP server. This server will listen for requests from our CLI during tests and return predefined JSON responses, mimicking the real API.

## 3. Implementation Steps

### Step 1: Create the Mock API Server

- **File:** `internal/gtasks/mock_server_test.go`
- **Details:**
  - Create a lightweight HTTP server using `httptest.NewServer`.
  - Implement a request router (`http.ServeMux`) to handle different API endpoints (e.g., `/api/v1/users/@me/lists`, `/api/v1/lists/{id}/tasks`).
  - For each endpoint, create a handler function that:
    - Checks the HTTP method (`GET`, `POST`, `DELETE`, etc.).
    - Constructs a realistic JSON response that matches the Google Tasks API format.
    - Writes the JSON response with the correct `Content-Type` header and status code.

### Step 2: Create a Test-Specific Client Constructor

- **File:** `internal/gtasks/client_test.go` (or similar)
- **Details:**
  - Create a new, un-exported function: `newTestClient(serverURL string)`.
  - This function will configure a `tasks.Service` to use our mock server's URL by using `option.WithEndpoint(serverURL)`.
  - It will also use a basic `http.Client` with `option.WithHTTPClient(&http.Client{})` to bypass the real OAuth2 flow.
  - It will return a `gtasks.Client` instance configured to talk to our mock server.

### Step 3: Update and "Un-skip" the Unit Tests

- **Files:** `internal/gtasks/tasklists_test.go`, `internal/gtasks/tasks_test.go`
- **Details:**
  - Remove the `t.Skip()` calls from the existing tests.
  - Each test function will be updated to follow this pattern:
    1.  **Start Mock Server:** `server := httptest.NewServer(...)`
    2.  **Defer Cleanup:** `defer server.Close()`
    3.  **Create Test Client:** `client := newTestClient(server.URL)`
    4.  **Capture Output:** Redirect `os.Stdout` to a buffer to capture the CLI's output.
    5.  **Execute Logic:** Call the function being tested (e.g., `client.ListTaskLists()`).
    6.  **Assert Results:**
        - Check for `nil` error.
        - Assert that the captured output contains the expected text from the mock server's response.

This plan will provide a robust and maintainable testing framework for our application.
