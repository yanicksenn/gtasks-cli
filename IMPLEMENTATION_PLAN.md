# Implementation Plan: Web-Based Authentication Flow

**Objective:** Replace the current local OAuth 2.0 flow in `gtasks-cli` with a modern, web-based authentication flow that relies on an external service.

---

### 1. Dependency Management

-   **Action:** Add the `github.com/pkg/browser` dependency to the project.
-   **Command:** `go get github.com/pkg/browser`
-   **Purpose:** To reliably open the user's default web browser to the authentication URL.

### 2. Create the Web Authentication Handler

-   **File:** `internal/auth/oauth_flow.go`
-   **Function:** `LoginViaWebFlow(ctx context.Context) (string, error)`
-   **Details:**
    1.  **Constants:** Define the `authServiceURL` (`https://oauth-hub-dev...`) and the local `redirectURI` (`http://localhost:8080/callback`).
    2.  **Local Server:**
        -   Start a local HTTP server on port `8080`.
        -   Create a `/callback` handler.
        -   Use a `chan` to signal the main function when the callback is received.
    3.  **Browser Interaction:**
        -   Construct the full login URL (e.g., `[authServiceURL]/login/gtasks`).
        -   Use `browser.OpenURL()` to open the URL.
        -   Provide a fallback message printing the URL to the console in case the browser cannot be opened automatically.
    4.  **Callback Handling:**
        -   The `/callback` handler will read the `access_token` and `refresh_token` from the request's query parameters.
        -   It will send these tokens back to the main `LoginViaWebFlow` function via the channel.
    5.  **User Info & Token Caching:**
        -   Once the tokens are received, use them to create an OAuth2 client.
        -   Call the Google User Info API to get the user's email address.
        -   Save the tokens to the `gtasks-token.json` cache, keyed by the user's email.
    6.  **Return Value:** Return the user's email on success or an error if any step fails.

### 3. Modify the `login` Command

-   **File:** `cmd/accounts.go`
-   **Target:** The `loginCmd.Run` function.
-   **Modifications:**
    1.  **Remove Old Logic:** Delete all code related to `auth.NewAuthenticator()` and the instructions for creating `credentials.json`.
    2.  **Call New Flow:** Replace the old logic with a single call to `auth.LoginViaWebFlow(context.Background())`.
    3.  **Update User Feedback:** Change the print statements to inform the user that their browser is opening and to provide success or failure messages based on the result of the web flow.

### 4. Refactor Existing Authentication Code

-   **File:** `internal/auth/auth.go`
-   **Actions:**
    1.  **Remove `NewAuthenticator`:** Delete the `NewAuthenticator` struct and the `NewAuthenticator()` constructor.
    2.  **Remove `NewClient`:** Delete the `NewClient` method, as its functionality is now handled by `LoginViaWebFlow`.
    3.  **Remove `loadCredentials`:** Delete the `Credentials` struct and the `loadCredentials` function, as the CLI no longer uses a local `credentials.json`.
    4.  **Preserve Core Utilities:** Keep the functions for token cache management (`loadTokenCache`, `save`, `saveAll`) and account management (`Logout`, `ListAccounts`), as they are still needed.

### 5. Update Tests

-   **Action:** Modify the tests for the `login` command to accommodate the new asynchronous, web-based flow.
-   **Strategy:**
    -   The test will need to run a mock HTTP server that simulates the external auth service.
    -   When the `login` command is executed in the test, it will open the browser to the mock server's URL.
    -   The mock server will then make a request to the CLI's local `/callback` endpoint, providing dummy tokens.
    -   The test will assert that the correct user email is returned and that the tokens are cached correctly.