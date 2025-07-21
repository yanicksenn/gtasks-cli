# End-to-End Tests

This directory contains the end-to-end (E2E) tests for the `gtasks` CLI.

## Table of Contents

- [Running the Tests](#running-the-tests)

---

## Running the Tests

To run the E2E tests, navigate to this directory and run:

```
go test -v
```

**Note:** The E2E tests require a valid `credentials.json` file in `~/.config/` and will perform a real OAuth2 flow, requiring user interaction.
