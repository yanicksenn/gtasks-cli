# OAuth Web Flow - Implementation Plan

This document outlines the plan for implementing a centralized, general-purpose OAuth flow using a Go-based web service hosted on Google Cloud Platform.

---

### Milestone 1: The General-Purpose Go OAuth Service

This service will be a standalone, reusable piece of infrastructure.

-   **1.1: Project Scaffolding & Configuration:**
    -   [ ] Create the directory `go/cmd/oauth-svc` and initialize a Go module.
    -   [ ] Create a `config.go` file. It will define structs for a `Config` object that contains a list of `Service` configurations.
    -   [ ] Use the `gopkg.in/yaml.v3` library to parse a `config.yaml` file at startup.
    -   [ ] **(Anti-Pattern Prevention):** The `config.yaml` will **not** contain secrets. It will only contain the *name of the environment variable* where the service should find the secret (e.g., `client_secret_env_var: "GTASKS_GOOGLE_CLIENT_SECRET"`).

-   **1.2: Core HTTP Server Logic:**
    -   [ ] In `main.go`, create a new `http.ServeMux` (router).
    -   [ ] The main handler will parse the URL path to extract the `{serviceName}` (e.g., `gtasks`). If the service name from the URL is not found in the loaded config, it will return a `404 Not Found` error.
    -   [ ] This handler will then delegate to specific `/login` and `/callback` handlers, passing the resolved service configuration to them.

-   **1.3: Implement `/login/{serviceName}` Endpoint:**
    -   [ ] This handler will receive the service-specific config.
    -   [ ] It will use the `golang.org/x/oauth2/google` package to create an `oauth2.Config` object using the `ClientID` and `ClientSecret` (read from the environment variable).
    -   [ ] It will then call the `AuthCodeURL` method on this config object to generate the correct, secure URL for Google's consent page.
    -   [ ] It will issue an HTTP `302 Found` redirect to send the user's browser to this URL.

-   **1.4: Implement `/callback/{serviceName}` Endpoint:**
    -   [ ] This handler will receive the `code` and `state` parameters from Google.
    -   [ ] It will perform the token exchange using the `oauth2.Config.Exchange()` method.
    -   [ ] Upon successful exchange, it will construct the redirect URL for the local CLI (e.g., `http://localhost:8080/callback`) and append the `access_token` and `refresh_token` as query parameters.
    -   [ ] It will then issue a `302 Found` redirect to the user's browser.

-   **1.5: User-Facing Error Handling:**
    -   [ ] **(Anti-Pattern Prevention):** Do not show raw errors or stack traces to the user.
    -   [ ] Create a simple `error.html` template.
    -   [ ] If any step in the `/callback` (like the token exchange) fails, the service will render this HTML template with a user-friendly message (e.g., "Authentication failed. Please try again.") and a proper HTTP status code (e.g., `500 Internal Server Error`).

-   **1.6: Dockerfile & Deployment (`README.md`):**
    -   [ ] Create a multi-stage `Dockerfile` using `golang:alpine` for building and `gcr.io/distroless/static-debian11` for the final, minimal image.
    -   [ ] The `Dockerfile` will copy the `config.yaml` into the image.
    -   [ ] Create a `README.md` for the `oauth-svc` with detailed GCP deployment instructions:
        -   How to use **Secret Manager** to securely store the `GTASKS_GOOGLE_CLIENT_SECRET`.
        -   How to build and push the image to **Artifact Registry**.
        -   The full `gcloud run deploy` command, showing how to mount the secret from Secret Manager as an environment variable.

---

### Milestone 2: CLI Modifications (`gtasks`)

-   **2.1: Streamline the `login` Command:**
    -   [ ] **(Anti-Pattern Prevention):** Remove the old, local OAuth flow entirely. The `login` command should have one job: to orchestrate the web flow.
-   **2.2: Implement Build-Time Variable for URL:**
    -   [ ] In a `config.go` or `main.go` file within `gtasks`, define `var AuthURL = "http://localhost:8081"`. The default value is for the developer's local instance of the auth service.
    -   [ ] Create a `BUILD.md` file in the project root explaining how to compile a production binary by injecting the production URL with `-ldflags`.
-   **2.3: Implement the Web Login Flow:**
    -   [ ] The `login` command will start a local server on a specific port (e.g., 8080).
    -   [ ] **(UX Improvement):** Use a library like `github.com/pkg/browser` to reliably open the user's default browser to `[AuthURL]/gtasks/login`.
    -   [ ] **(Graceful Fallback):** If opening the browser fails (e.g., in an SSH session), the CLI will not crash. Instead, it will print the URL to the console and instruct the user to copy and paste it into a browser manually.
    -   [ ] The local `/callback` handler will be set up to wait for the request containing the tokens.
    -   [ ] **(Robust Shutdown):** Use a Go channel to signal the main `login` command that the tokens have been received. This is a clean way to block the command until the flow is complete and then allow the local server to shut down gracefully.
-   **2.4: Update Tests:**
    -   [ ] The tests for the login flow will be updated to reflect this new, streamlined process.

---

### Milestone 3: Documentation and Finalization

-   **3.1: Update `README.md`:**
    -   [ ] Update the main `gtasks` `README.md` to reflect the new, simpler login process.
-   **3.2: Update `DESIGN.md`:**
    -   [ ] Update the "Authentication Flow" section to describe the new architecture.
-   **3.3: Update `WORKLOG.md`:**
    -   [ ] Add a detailed entry for this major feature.