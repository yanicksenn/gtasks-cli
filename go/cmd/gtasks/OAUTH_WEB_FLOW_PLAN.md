# OAuth Web Flow - Implementation Plan

This document outlines the plan for implementing a centralized OAuth flow using a Go-based web service hosted on Google Cloud Platform.

---

### Milestone 1: The Go Authentication Web Service

-   **1.1: Project Scaffolding:**
    -   [ ] Create a new directory: `go/cmd/gtasks-auth-svc`.
    -   [ ] Initialize a new Go module within it.
    -   [ ] Use Go's standard `net/http` package for the web server.
-   **1.2: Configuration:**
    -   [ ] The service will be configured entirely via environment variables (`GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`). Create a simple `config.go` file to load these.
-   **1.3: Implement `/login` Endpoint:**
    -   [ ] This endpoint will construct the Google OAuth2 URL.
    -   [ ] It will include a `redirect_uri` parameter that points back to our service's own `/callback` endpoint.
    -   [ ] It will then issue an HTTP `302 Found` redirect to send the user's browser to Google's consent page.
-   **1.4: Implement `/callback` Endpoint:**
    -   [ ] This endpoint will receive the `code` from Google.
    -   [ ] It will use Go's `golang.org/x/oauth2` library to securely exchange the `code` for an `access_token` and `refresh_token` by making a server-to-server call to Google.
    -   [ ] It will then construct a new URL for the CLI's local server (e.g., `http://localhost:8080/callback?access_token=...&refresh_token=...`) and issue another `302 Found` redirect to send the user's browser back to the CLI.
-   **1.5: Create a `Dockerfile`:**
    -   [ ] Create a multi-stage `Dockerfile` to build a small, optimized, and secure container image for our service.
-   **1.6: GCP Deployment Documentation:**
    -   [ ] Create a `README.md` for the service with specific, copy-pasteable instructions for deploying it to **Google Cloud Run**. The documentation will cover:
        -   How to build the Docker image and push it to Google Artifact Registry.
        -   The `gcloud run deploy` command needed to deploy the service.
        -   How to securely set the `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` as secrets in Cloud Run.

---

### Milestone 2: CLI Modifications

-   **2.1: Add `--web` Flag:**
    -   [ ] Add a `--web` flag to the `gtasks login` command.
-   **2.2: Implement Web Login Logic:**
    -   [ ] The `Run` function for `login --web` will:
        -   Start a temporary local HTTP server on `localhost:8080`.
        -   This server will have a `/callback` endpoint that expects `access_token` and `refresh_token` as query parameters.
        -   It will open the user's browser to the deployed Cloud Run service's `/login` URL.
        -   Upon receiving the callback, it will save the tokens to the local cache, print a success message, and shut down the local server.
-   **2.3: Update Tests:**
    -   [ ] Add a new test case for the web login flow. This test will use the `httptest` package to mock the interaction with our new auth web service, simulating the redirect flow without making real network calls.

---

### Milestone 3: Documentation and Finalization

-   **3.1: Update `README.md`:**
    -   [ ] Update the main `README.md` to explain the new, simpler web-based login as the primary method for users.
-   **3.2: Update `DESIGN.md`:**
    -   [ ] Update the "Authentication Flow" section to describe the new architecture involving the Cloud Run service.
-   **3.3: Update `WORKLOG.md`:**
    -   [ ] Add an entry summarizing this major feature addition.
