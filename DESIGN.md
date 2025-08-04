# gtasks CLI - Software Design

**Related Documents:**
- [Requirements (`README.md`)](./README.md)
- [Implementation Plan (`IMPLEMENTATION_PLAN.md`)](./IMPLEMENTATION_PLAN.md)

---

## Table of Contents

- [1. Overall Architecture](#1-overall-architecture)
- [2. Directory Structure](#2-directory-structure)
- [3. Authentication Flow (OAuth 2.0)](#3-authentication-flow-oauth-20)
- [4. Configuration Management](#4-configuration-management)
- [5. API and Client Architecture](#5-api-and-client-architecture)
  - [The Client Interface (Service Facade)](#the-client-interface-service-facade)
  - [Online and Offline Clients](#online-and-offline-clients)
- [6. Offline Mode](#6-offline-mode)
  - [Persistent Offline Store](#persistent-offline-store)
- [7. Testing Strategy](#7-testing-strategy)](#7-testing-strategy)
  - [Offline Tests (Unit/Integration)](#offline-tests-unitintegration)
  - [End-to-End (E2E) Tests](#end-to-end-e2e-tests)

---

This document outlines the high-level software design for the `gtasks` CLI tool.

## 1. Overall Architecture

The application is built in Go and structured around the **Cobra** library. It uses a clean architecture that separates the command-line interface (the `cmd` package) from the business logic (the `internal/gtasks` package).

## 2. Directory Structure

The project follows standard Go conventions:
```
.
├── cmd/          # Cobra command definitions
├── internal/
│   ├── auth/     # Google OAuth2 authentication
│   ├── config/   # Configuration handling
│   └── gtasks/   # Core business logic and API clients
├── e2e/          # End-to-end tests
└── main.go       # Main application entry point
```

## 3. Authentication Flow (OAuth 2.0)

Authentication is handled via a standard OAuth 2.0 flow for desktop applications. Credentials are cached locally to provide a seamless experience after the initial login.

## 4. Configuration Management

A simple configuration file (`~/.config/gtasks/config.yaml`) stores the active user account.

## 5. API and Client Architecture

### The Client Interface (Service Facade)

To support both online and offline modes, the application uses a `Client` interface. This interface acts as a "service facade," defining all the high-level operations the CLI can perform (e.g., `ListTaskLists`, `CreateTask`).

The command layer of the application interacts only with this interface, making it completely unaware of whether it is operating on the real Google Tasks API or a local offline store.

### Online and Offline Clients

There are two concrete implementations of the `Client` interface:
- **`onlineClient`:** This client contains the official Google Tasks API service. Its methods make live HTTP requests to Google's servers.
- **`offlineClient`:** This client contains a persistent, file-based store. Its methods read from and write to a local JSON file, performing the same operations as the `onlineClient` but without any network access.

A `NewClient` factory function is responsible for creating and returning the correct client implementation based on whether the user has provided the `--offline` flag.

## 6. Offline Mode

### Persistent Offline Store

The offline mode is powered by a persistent `offlineStore`. This is a thread-safe, file-based database (a simple JSON file at `~/.config/gtasks/offline.json`) that stores all task lists and tasks created while offline.

## 7. Testing Strategy

The project employs a two-tiered testing strategy.

### Offline Tests (Unit/Integration)

The primary testing method is a suite of offline tests that run against a stateful, in-memory mock of the Google Tasks API. This allows for fast, deterministic, and authentication-free testing of the entire business logic.

### End-to-End (E2E) Tests

A small suite of E2E tests compiles and runs the actual CLI binary to verify basic functionality like the help command.
