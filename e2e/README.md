# End-to-End Tests

This directory contains the end-to-end (E2E) tests for the `gtasks` CLI.

## Table of Contents

- [Running the Tests](#running-the-tests)

---

## Running the Tests

To run the E2E tests, navigate to the root of the project and run:

```
go test ./...
```

**Note:** The E2E tests will perform a real OAuth2 flow, but the credentials are hardcoded and the login/logout process is handled automatically by the test suite.
